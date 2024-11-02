package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	chat_actions "github.com/R1kkass/GoCloudGRPC/actions/chat"
	"github.com/R1kkass/GoCloudGRPC/db"
	"github.com/R1kkass/GoCloudGRPC/helpers"
	Model "github.com/R1kkass/GoCloudGRPC/models"
	"github.com/R1kkass/GoCloudGRPC/proto/chat"
	"github.com/R1kkass/GoCloudGRPC/structs"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func CreateChat(ctx context.Context, in *chat.CreateRequestChat) (*chat.CreateResponseChat, error) {

	var chats Model.Chat

	user, err := helpers.GetUserFormMd(ctx)

	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Пользователь не найден")
	}

	if err := chat_actions.CheckChatExist(ctx, in); err != nil {
		return &chat.CreateResponseChat{
			Message: "Ошибка",
		}, status.Error(codes.AlreadyExists, "Чат уже существует")
	}

	var p string
	var g int64

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(
			&chats,
		)
		var chatUsers Model.ChatUser

		if result.Error != nil {
			return errors.New("чат не создан")
		}

		result = tx.Model(&chatUsers).Create(
			&Model.ChatUser{
				ChatID: int(chats.ID),
				UserRelation: Model.UserRelation{
					UserID: int(user.ID),
				},
				SubmitCreate: true,
			},
		).Create(&Model.ChatUser{
			ChatID: int(chats.ID),
			UserRelation: Model.UserRelation{
				UserID: int(in.GetOtherId()),
			},
			SubmitCreate: false,
		})
		if result.Error != nil {
			return errors.New("чат не создан")
		}

		p, g = chat_actions.SendFirstParams(&chats)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	objectMessage := map[string]any{
		"title": "Новый запрос на отправку сообщений",
		"type":  "New_ChatRequest",
	}

	go chat_actions.NotificationChatCreate(int(in.GetOtherId()), objectMessage)

	return &chat.CreateResponseChat{
		Message: "Чат создан",
		Keys: &chat.Keys{
			P: p,
			G: g,
		},
		ChatId: uint32(chats.ID),
	}, nil
}

func StreamGetChat(in *chat.Empty, requestStream chat.ChatGreeter_StreamGetChatServer) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error get chats: ", r)
		}
	}()
	var chats []*chat.ChatUsersCount

	ctx := requestStream.Context()
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return status.Error(codes.Unauthenticated, "Токен не найден")
	}

	jwtToken, ok := md["authorization"]

	if !ok {
		return status.Error(codes.Unauthenticated, "Токен не найден")
	}

	user, err := helpers.GetUser(jwtToken)
	channel := make(chan map[string]any)

	if err != nil {
		return status.Error(codes.Unauthenticated, "Пользователь не найден")
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
		key := strconv.Itoa(int(user.ID)) + "_messages"

		for {
			res := db.ConnectRedisNotificationDB.Subscribe(ctx, key)
			var jsonDecodeMsg map[string]any
			message, err := res.ReceiveMessage(ctx)
			json.Unmarshal([]byte(message.Payload), &jsonDecodeMsg)

			if err != nil {
				log.Println("can't get chats")
				return
			}

			channel <- jsonDecodeMsg
		}

	}()

	db.DB.Model(&Model.ChatUser{}).Select(`chat_users.*, COALESCE(messages.created_at,'2022-10-19 15:23:53.252567+00') as create_at_message, count(un_readed_messages.id) as un_readed_messages_count`).
		Joins("LEFT JOIN un_readed_messages ON un_readed_messages.chat_id = chat_users.chat_id AND un_readed_messages.user_id = ?", user.ID).
		Joins("LEFT JOIN (SELECT * FROM (SELECT distinct on(chat_id) chat_id, created_at FROM messages ORDER BY chat_id, created_at DESC) t ORDER BY created_at DESC) AS messages ON messages.chat_id = chat_users.chat_id", user.ID).
		Preload("User").Preload("Chat").Preload("Chat.ChatUsers.User").
		Preload("Chat.Message", func(db *gorm.DB) *gorm.DB {
			return db.Order("messages.id ASC")
		}).
		Where("chat_users.user_id = ? AND submit_create = ?", user.ID, true).
		Group("chat_users.id, un_readed_messages.chat_id, messages.created_at").
		Order("create_at_message DESC, chat_users.created_at DESC").
		Find(&chats)

	err = requestStream.Send(&chat.StreamGetResponseChat{
		Chats: chats,
	})
	if err != nil {
		log.Println("error while sending chats:", err)
		return err
	}

	for {
		<-channel

		var chats []*chat.ChatUsersCount

		db.DB.Model(&Model.ChatUser{}).Select(`chat_users.*, COALESCE(messages.created_at,'2022-10-19 15:23:53.252567+00') as create_at_message, count(un_readed_messages.id) as un_readed_messages_count`).
			Joins("LEFT JOIN un_readed_messages ON un_readed_messages.chat_id = chat_users.chat_id AND un_readed_messages.user_id = ?", user.ID).
			Joins("LEFT JOIN (SELECT * FROM (SELECT distinct on(chat_id) chat_id, created_at FROM messages ORDER BY chat_id, created_at DESC) t ORDER BY created_at DESC) AS messages ON messages.chat_id = chat_users.chat_id", user.ID).
			Preload("User").Preload("Chat").Preload("Chat.ChatUsers.User").
			Preload("Chat.Message", func(db *gorm.DB) *gorm.DB {
				return db.Order("messages.id ASC")
			}).
			Where("chat_users.user_id = ? AND submit_create = ?", user.ID, true).
			Group("chat_users.id, un_readed_messages.chat_id, messages.created_at").
			Order("create_at_message DESC, chat_users.created_at DESC").
			Find(&chats)

		err := requestStream.Send(&chat.StreamGetResponseChat{
			Chats: chats,
		})
		if err != nil {
			log.Println("error while sending chats:", err)
			return err
		}

	}
}

func GetUnSuccessChats(ctx context.Context, in *chat.Empty) (*chat.GetUnSuccessChatsResponse, error) {
	var chats []*chat.ChatUsers
	user, err := helpers.GetUserFormMd(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Пользователь не найден")
	}

	db.DB.Model(&Model.ChatUser{}).
		Preload("User").Preload("Chat").Preload("Chat.ChatUsers.User").
		Where("chat_users.user_id = ? AND submit_create = ?", user.ID, false).
		Find(&chats)

	return &chat.GetUnSuccessChatsResponse{
		Chats: chats,
	}, nil
}

func CreateSecondaryKey(ctx context.Context, in *chat.CreateSecondaryKeyRequest) (*chat.CreateSecondaryKeyResponse, error) {
	user, err := helpers.GetUserFormMd(ctx)

	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Пользователь не найден")
	}
	_, err = chat_actions.CheckSecondaryKey(user.ID, in.GetChatId())

	if err != nil {
		return nil, status.Error(codes.AlreadyExists, "Ключ уже создан")
	}

	db.DB.Create(&Model.KeysSecondary{
		UserID: user.ID,
		ChatID: uint(in.GetChatId()),
		Key:    in.GetKey(),
	})

	return &chat.CreateSecondaryKeyResponse{
		Message: "Успешно",
	}, nil
}

func GetSecondaryKey(ctx context.Context, in *chat.GetSecondaryKeyRequest) (*chat.GetSecondaryKeyResponse, error) {

	user, err := helpers.GetUserFormMd(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Пользователь не найден")
	}

	err, _ = chat_actions.CheckChat(in.GetChatId(), user.ID)

	if err != nil {
		return nil, status.Error(codes.NotFound, "Чат не найден")
	}

	key, err := chat_actions.GetSecondaryKey(user.ID, in.GetChatId())

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	keys, err := chat_actions.GetPublicKey(in.GetChatId())

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &chat.GetSecondaryKeyResponse{
		Key: key.Key,
		P:   keys.P,
	}, nil
}

func GetPublicKey(ctx context.Context, in *chat.GetPublicKeyRequest) (*chat.GetPublicKeyResponse, error) {

	if err := helpers.CheckChat(ctx, in.GetChatId()); err != nil {
		return nil, status.Error(codes.NotFound, "Чат не найден")
	}

	keys, err := chat_actions.GetPublicKey(in.GetChatId())

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &chat.GetPublicKeyResponse{
		G: keys.G,
		P: keys.P,
	}, nil
}

func AcceptChat(ctx context.Context, in *chat.AcceptChatRequest) (*chat.AcceptChatResponse, error) {

	if err := helpers.CheckChat(ctx, in.GetChatId()); err != nil {
		return nil, status.Error(codes.NotFound, "Чат не найден")
	}

	user, err := helpers.GetUserFormMd(ctx)

	if err != nil {
		return nil, status.Error(codes.Aborted, "Пользователь не найден")
	}

	db.DB.Model(&Model.ChatUser{}).Where("user_id = ? AND chat_id = ?", user.ID, in.GetChatId()).Update("submit_create", true)

	return &chat.AcceptChatResponse{
		Message: "Успех",
	}, nil
}

func DissalowChat(ctx context.Context, in *chat.DissalowChatRequest) (*chat.DissalowChatResponse, error) {
	var chatUser Model.ChatUser
	if err := helpers.CheckChat(ctx, in.GetChatId()); err != nil {
		return nil, status.Error(codes.NotFound, "Чат не найден")
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("error ", r)
		}
	}()
	user, err := helpers.GetUserFormMd(ctx)

	if err != nil {
		return nil, status.Error(codes.NotFound, "Пользователь не найден")
	}

	result := db.DB.Model(&Model.ChatUser{}).Where("chat_id = ? AND user_id = ?", in.GetChatId(), user.ID).First(&chatUser)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, status.Error(codes.Aborted, "Чат не найден")
	}

	db.DB.Transaction(func(tx *gorm.DB) error {
		result = tx.Where("chat_id = ?", chatUser.ChatID).Unscoped().Delete(&Model.ChatUser{})

		if result.Error != nil {
			return errors.New("ошибка")
		}

		result := tx.Where("id=?", chatUser.ChatID).Unscoped().Delete(&Model.Chat{})
		if result.Error != nil {
			return errors.New("ошибка")
		}

		return nil
	})

	return &chat.DissalowChatResponse{
		Message: "Успех",
	}, nil
}

func GetMessages(ctx context.Context, in *chat.GetMessagesRequest) (*chat.GetMessagesResponse, error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error get messages: ", r)
		}
	}()

	var message []*chat.Message

	if err := helpers.CheckChat(ctx, in.GetChatId()); err != nil {
		return nil, status.Error(codes.NotFound, "Чат не найден")
	}

	user, err := helpers.GetUserFormMd(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Пользователь не найден")
	}
	var page = int64(in.GetPage())
	var count int64 = 0
	count, err = chat_actions.GetCountNotReadedMessages(int(in.GetChatId()), int(user.ID))

	if err != nil {
		return nil, status.Error(codes.Unknown, "Не удалось получить сообщения")
	}

	if in.GetInit() {
		page = count / 10
		if page != 0 {
			page += 1
		}
	}

	r := db.DB.Model(&Model.Message{}).Preload("User").
		Select("messages.*, SUM(CASE WHEN un_readed_messages.id IS NULL THEN 0 ELSE 1 END) AS un_readed_message").
		Joins("LEFT JOIN un_readed_messages ON un_readed_messages.message_id = messages.id AND un_readed_messages.user_id = ?", user.ID).
		Where("messages.chat_id = ?", in.GetChatId()).
		Group("messages.id").
		Order("id DESC").Offset(10 * int(page)).
		Limit(10).
		Find(&message)

	if r.Error != nil {
		return nil, status.Error(codes.Unknown, "Не удалось получить сообщения")
	}

	return &chat.GetMessagesResponse{
		Messages:     message,
		Page:         int32(page),
		CountNotRead: int32(count),
	}, nil
}

func StreamGetMessagesGeneral(in *chat.Empty, responseStream chat.ChatGreeter_StreamGetMessagesGeneralServer) error {
	ctx := responseStream.Context()
	user, err := helpers.GetUserFormMd(ctx)
	channel := make(chan bool)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error get general messages: ", r)
		}
	}()

	if err != nil {
		return status.Error(codes.Unauthenticated, "Пользователь не найден")
	}

	var count int64
	r := db.DB.Model(&Model.UnReadedMessage{}).Where("user_id = ?", user.ID).Distinct("chat_id").Count(&count)

	if r.Error != nil {
		log.Println("error while sending count messages:", err)
		return err
	}

	err = responseStream.Send(&chat.StreamGetMessagesGeneralResponse{
		Count: int32(count),
	})

	if err != nil {
		log.Println("error while sending count messages:", err)
		return err
	}

	go chat_actions.NotificationObserver(ctx, int(user.ID), &channel)

	for {
		<-channel
		var count int64
		r := db.DB.Model(&Model.UnReadedMessage{}).Where("user_id = ?", user.ID).Distinct("chat_id").Count(&count)
		if r.Error != nil {
			log.Println("error while sending count messages:", err)
			return err
		}
		err := responseStream.Send(&chat.StreamGetMessagesGeneralResponse{
			Count: int32(count),
		})

		if err != nil {
			log.Println("error while sending count messages:", err)
			return err
		}
	}
}

func StreamGetMessages(stream chat.ChatGreeter_StreamGetMessagesServer, conns map[string]structs.DataStreamConnect, chatId int, userId int, channel *chan *chat.StreamGetMessagesResponse, token string) error {

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
				ChatID: chatId,
			}

			r := db.DB.Preload("User").
				Create(&message).
				Select("messages.*, true as un_readed_message").
				Where("messages.id = ?", message.ID).First(&messageResponse)

			if r.Error != nil || r.RowsAffected == 0 {
				CloseConnect(conns, token)
				return status.Error(codes.Unknown, "неизвестная ошибка")
			}

			go chat_actions.NotificationMessageCreate(chatId, msg.GetMessage(), userId, conns, message.ID)

			*channel <- &chat.StreamGetMessagesResponse{
				Message: messageResponse,
				Type:    msg.GetType(),
			}
		}

		if msg.GetType() == chat.TypeMessage_READ_MESSAGE {
			message, err := chat_actions.RemoveUnReadedMessage(int(msg.GetMessageId()), userId, chatId)

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
					db.ConnectRedisDB.Publish(context.TODO(), strconv.Itoa(userId)+"_messages_general", "true")
				}
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
