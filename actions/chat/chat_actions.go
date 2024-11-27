package chat_actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/R1kkass/GoCloudGRPC/db"
	"github.com/R1kkass/GoCloudGRPC/helpers"
	Model "github.com/R1kkass/GoCloudGRPC/models"
	"github.com/R1kkass/GoCloudGRPC/proto/chat"
	"github.com/R1kkass/GoCloudGRPC/structs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SendFirstParams(chat *Model.Chat) (string, int64) {
	p, g := helpers.GeneratePubKeys()

	db.DB.Create(&Model.Keys{
		ChatID: chat.ID,
		P:      p.String(),
		G:      g,
	})

	return p.String(), g
}

func SendSecondaryParams(user *Model.User, chat *Model.Chat, key string) error {
	result := db.DB.Create(&Model.KeysSecondary{
		UserID: user.ID,
		ChatID: chat.ID,
		Key:    key,
	})

	if result.Error != nil {
		return errors.New("ключ не создан")
	}

	return nil
}

func GetSecondaryKey(user_id uint, chat_id uint32) (*Model.KeysSecondary, error) {
	var keys *Model.KeysSecondary
	result := db.DB.Model(&Model.KeysSecondary{}).Where("chat_id = ? AND user_id <> ?", chat_id, user_id).First(&keys)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil, errors.New("ключ не найден")
	}

	return keys, nil
}

func CheckSecondaryKey(user_id uint, chat_id uint32) (*Model.KeysSecondary, error) {
	var keys *Model.KeysSecondary
	result := db.DB.Model(&Model.KeysSecondary{}).Where("chat_id = ? AND user_id = ?", chat_id, user_id).First(&keys)

	fmt.Println(user_id, chat_id)

	if result.RowsAffected != 0 {
		return nil, errors.New("ключ уже существует")
	}

	return keys, nil
}

func CheckChat(chat_id uint32, user_id uint) (error, *Model.ChatUser) {
	var chat *Model.ChatUser

	result := db.DB.Model(&Model.ChatUser{}).Where("chat_id = ? AND user_id = ?", chat_id, user_id).First(&chat)

	if result.RowsAffected == 0 {
		return errors.New("чат не найден"), nil
	}

	return nil, chat
}

func CheckChatByMessageId(message_id uint32, user_id uint) (*Model.ChatUser, error) {
	var message *Model.Message

	r := db.DB.Model(&Model.Message{}).Where("id = ? AND user_id = ?", message_id, user_id).First(&message)

	if r.Error != nil || r.RowsAffected == 0 {
		return nil, errors.New("сообщение не найдено")
	}

	var chat *Model.ChatUser

	result := db.DB.Model(&Model.ChatUser{}).Where("chat_id = ? AND user_id = ?", message.ChatID, user_id).First(&chat)

	if result.RowsAffected == 0 {
		return nil, errors.New("чат не найден")
	}

	return chat, nil
}

func GetPublicKey(chatId uint32) (*Model.Keys, error) {

	var keys *Model.Keys
	result := db.DB.Model(&Model.Keys{}).Where("chat_id = ?", chatId).First(&keys)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, errors.New("ключ не найден")
	}

	return keys, nil
}

func CheckChatExist(ctx context.Context, in *chat.CreateRequestChat) error {
	user, _ := helpers.GetUserFormMd(ctx)

	var usersChat Model.ChatUser

	result := db.DB.Raw(`SELECT count(*), chat_id from chat_users WHERE chat_id in (SELECT chat_id FROM chat_users Where user_id = ? INTERSECT SELECT chat_id FROM chat_users Where user_id = ?) GROUP BY chat_id`, user.ID, in.GetOtherId()).Scan(&usersChat)

	if result.RowsAffected != 0 {
		return errors.New("чат уже существует")
	}

	return nil
}

func NotificationChatCreate(userId int, objectMessage map[string]any) {
	obj, err := json.Marshal(objectMessage)

	if err != nil {
		fmt.Println("Error send notification: ", err)
		return
	}

	r := db.ConnectRedisDB.Publish(context.TODO(), strconv.Itoa(userId)+"_notification", obj)
	if r.Err() != nil {
		fmt.Println("Error send notification: ", r.Err().Error())
	}
}

func NotificationMessageCreate(chatId uint, message string, userId uint, conns map[string]structs.DataStreamConnect, messageId uint) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error in NotificationMessageCreate: ", r)
		}
	}()

	mapMessage := map[string]any{
		"description": message,
		"title":       "Новое сообщение",
		"type":        "New_Message",
		"options": map[string]string{
			"chat_id": strconv.Itoa(int(chatId)),
		},
	}
	var users []*Model.ChatUser

	r := db.DB.Model(&Model.ChatUser{}).Where("chat_id = ?", chatId).Find(&users)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Error in NotificationMessageCreate: ", r)
			}
		}()

		var connectedUsers = make(map[uint]bool)
		for _, connectedUser := range conns {
			connectedUsers[connectedUser.UserID] = true
		}
		objectMessage, _ := json.Marshal(mapMessage)

		for _, user := range users {
			_, ok := connectedUsers[user.UserID]
			if user.UserID != userId && !ok {
				db.ConnectRedisDB.Publish(context.TODO(), strconv.Itoa(int(user.UserID))+"_notification", objectMessage)
			}
			if user.UserID != userId {
				message := &Model.UnReadedMessage{

					ChatRelations: Model.ChatRelations{
						ChatID: chatId,
					},
					UserRelation: Model.UserRelation{
						UserID: user.UserID,
					},
					MessageRelations: Model.MessageRelations{
						MessageID: messageId,
					},
				}
				db.DB.Create(&message)
				go db.ConnectRedisDB.Publish(context.TODO(), strconv.Itoa(int(user.UserID))+"_messages", "true")
				go db.ConnectRedisDB.Publish(context.TODO(), strconv.Itoa(int(user.UserID))+"_messages_general", "true")
			}
		}
	}()

	if r.RowsAffected == 0 || r.Error != nil {
		return
	}
}

func NotificationObserver(ctx context.Context, userId int, channel *chan bool) {
	key := strconv.Itoa(userId) + "_messages_general"
	res := db.ConnectRedisNotificationDB.Subscribe(ctx, key)

	for {
		_, err := res.ReceiveMessage(ctx)

		if err != nil {
			log.Println("Can not create subscribe")
			break
		}

		*channel <- true
	}
}

func RemoveUnReadedMessage(messageId uint, userId uint, chatId uint) (*chat.Message, error) {
	r := db.DB.Unscoped().Where("user_id = ? AND message_id = ? AND chat_id = ?", userId, messageId, chatId).Delete(&Model.UnReadedMessage{})

	if r.Error != nil {
		fmt.Println("error in RemoveUnReadedMessage: ", r.Error)
		return nil, errors.New("не удалось изменить статус сообщения")
	}
	go db.ConnectRedisDB.Publish(context.TODO(), strconv.Itoa(int(userId))+"_messages", "true")

	var message *chat.Message

	r = db.DB.Model(&Model.Message{}).Where("messages.id = ? AND messages.chat_id = ? AND messages.user_id != ? AND status_message = 'success'", messageId, chatId, userId).
		Select("messages.*, SUM(CASE WHEN un_readed_messages.id IS NULL THEN 0 ELSE 1 END) AS un_readed_message").
		Preload("User").
		Preload("Chatfiles").
		Joins("LEFT JOIN un_readed_messages ON un_readed_messages.message_id = messages.id AND un_readed_messages.user_id = ?", userId).
		Group("messages.id").
		First(&message)

	if r.Error != nil || r.RowsAffected == 0 {
		fmt.Println("error in RemoveUnReadedMessage: ", r.Error)
		return nil, errors.New("не удалось изменить статус сообщения")
	}

	if !message.StatusRead {
		r = db.DB.Model(&Model.Message{}).Where("id = ?", message.Id).Update("status_read", true)
		message.StatusRead = true
	} else {
		return nil, nil
	}

	if r.Error != nil {
		fmt.Println("error in RemoveUnReadedMessage: ", r.Error)
		return nil, errors.New("не удалось изменить статус сообщения")
	}

	return message, nil
}

func GetCountNotReadedMessages(chatId int, userId int) (int64, error) {
	var count int64
	r := db.DB.Model(&Model.UnReadedMessage{}).Where("chat_id = ? AND user_id = ?", chatId, userId).Count(&count)

	if r.Error != nil || r.RowsAffected == 0 {
		return 0, r.Error
	}

	return count, nil
}

func StreamGetMessages(stream chat.ChatGreeter_StreamGetMessagesServer, conns map[string]structs.DataStreamConnect, chatId uint, userId uint, channel *chan *chat.StreamGetMessagesResponse, token string) error {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error chat_actions.StreamGetMessages: ", r)
		}
	}()

	for {
		msg, err := stream.Recv()

		if err != nil {
			CloseConnect(conns, token)
			return status.Error(codes.Unknown, "неизвестная ошибка")
		}
		if msg.GetType() == chat.TypeMessage_SEND_MESSAGE {

			if len(msg.GetMessage()) < 100 {
				CloseConnect(conns, token)
				return status.Error(codes.Unknown, "неизвестная ошибка")
			}
			var messageResponse *chat.Message

			message := &Model.Message{
				Text: string(msg.GetMessage()),
				UserRelation: Model.UserRelation{
					UserID: userId,
				},
				ChatRelations: Model.ChatRelations{
					ChatID: chatId,
				},
				TypeMessage: Model.TextMessage,
			}

			r := db.DB.Preload("User").
				Create(&message).
				Select("messages.*, true as un_readed_message").
				Preload("ChatFiles").
				Where("messages.id = ?", message.ID).First(&messageResponse)

			if r.Error != nil || r.RowsAffected == 0 {
				CloseConnect(conns, token)
				return status.Error(codes.Unknown, "неизвестная ошибка")
			}

			go NotificationMessageCreate(chatId, msg.GetMessage(), userId, conns, message.ID)

			*channel <- &chat.StreamGetMessagesResponse{
				Message: messageResponse,
				Type:    msg.GetType(),
			}
		}

		if msg.GetType() == chat.TypeMessage_READ_MESSAGE {
			message, err := RemoveUnReadedMessage(uint(msg.GetMessageId()), userId, chatId)

			if err != nil {
				CloseConnect(conns, token)
				return status.Error(codes.Unknown, "неизвестная ошибка")
			}

			if message != nil {
				*channel <- &chat.StreamGetMessagesResponse{
					Message: message,
					Type:    msg.GetType(),
				}
			}
			var count int64
			r := db.DB.Model(&Model.UnReadedMessage{}).Where("chat_id = ? AND user_id = ?", chatId, userId).Count(&count)

			if r.Error == nil {
				if count == 0 {
					db.ConnectRedisDB.Publish(context.TODO(), strconv.Itoa(int(userId))+"_messages_general", "true")
				}
			}
		}

		if msg.GetType() == chat.TypeMessage_UPLOAD_MESSAGE {
			var message *chat.Message

			r := db.DB.Model(&Model.Message{}).
				Where("id = ?", msg.GetMessageId()).Update("status_message", Model.Success)

			if r.Error != nil || r.RowsAffected == 0 {
				db.DB.Unscoped().Where("id = ?", msg.GetMessageId()).Delete(&Model.Message{})
				CloseConnect(conns, token)
				return status.Error(codes.Unknown, "неизвестная ошибка")
			}
			db.DB.Model(&Model.Message{}).
				Preload("User").
				Preload("ChatFiles").
				Select("messages.*, true as un_readed_message").
				Where("id = ?", msg.GetMessageId()).
				First(&message)
			fmt.Println(message.UserId)
			go NotificationMessageCreate(chatId, message.Text, userId, conns, uint(message.Id))
			*channel <- &chat.StreamGetMessagesResponse{
				Message: message,
				Type:    chat.TypeMessage_SEND_MESSAGE,
			}
		}
	}
}

func CloseConnect(conns map[string]structs.DataStreamConnect, token string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error close connection: ", r)
		}
	}()
	delete(conns, token)
}
