package chat_actions

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/R1kkass/GoCloudGRPC/db"
	Model "github.com/R1kkass/GoCloudGRPC/models"
	"github.com/R1kkass/GoCloudGRPC/proto/chat"
)

func RemoveUnReadedMessageChat(userId uint, chatId uint) ([]*chat.Message, error) {
	r := db.DB.Unscoped().Where("user_id = ? AND chat_id = ?", userId, chatId).Delete(&Model.UnReadedMessage{})

	if r.Error != nil {
		fmt.Println("error in RemoveUnReadedMessageChat: ", r.Error)
		return nil, errors.New("не удалось изменить статус сообщения")
	}
	go db.ConnectRedisDB.Publish(context.TODO(), strconv.Itoa(int(userId))+"_messages", "true")

	var messages []*chat.Message

	r = db.DB.Model(&Model.Message{}).Where("messages.chat_id = ? AND messages.user_id != ? AND status_message = 'success'", chatId, userId).
		Select("messages.*, SUM(CASE WHEN un_readed_messages.id IS NULL THEN 0 ELSE 1 END) AS un_readed_message").
		Preload("User").
		Preload("ChatFiles").
		Joins("LEFT JOIN un_readed_messages ON un_readed_messages.message_id = messages.id AND un_readed_messages.user_id = ?", userId).
		Group("messages.id").
		Find(&messages)

	if r.Error != nil || r.RowsAffected == 0 {
		fmt.Println("error in RemoveUnReadedMessage: ", r.Error)
		return nil, errors.New("не удалось изменить статус сообщения")
	}

	for _, message := range messages {
		if !message.StatusRead {
			r = db.DB.Model(&Model.Message{}).Where("id = ?", message.Id).Update("status_read", true)
			message.StatusRead = true
		} else {
			return nil, nil
		}
	}

	if r.Error != nil {
		fmt.Println("error in RemoveUnReadedMessage: ", r.Error)
		return nil, errors.New("не удалось изменить статус сообщения")
	}

	return messages, nil
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
		Preload("ChatFiles").
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