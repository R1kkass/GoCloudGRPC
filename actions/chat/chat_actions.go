package chat_actions

import (
	"context"
	"errors"
	"fmt"
	"mypackages/db"
	"mypackages/helpers"
	Model "mypackages/models"
	"mypackages/proto/chat"
)

func SendFirstParams(chat *Model.Chat) (string, int64) {
	p,g := helpers.GeneratePubKeys()

	db.DB.Create(&Model.Keys{
		ChatID: chat.ID,
		P: p.String(),
		G: g,
	})

	return p.String(), g
}

func SendSecondaryParams(user *Model.User, chat *Model.Chat, key string) error {
	result := db.DB.Create(&Model.KeysSecondary{
		UserID: user.ID,
		ChatID: chat.ID,
		Key: key,
	})

	if result.Error != nil {
		return errors.New("ключ не создан")
	}

	return nil
}

func GetSecondaryKey(user_id uint, chat_id uint32) (*Model.KeysSecondary, error){
	var keys *Model.KeysSecondary 
	result := db.DB.Model(&Model.KeysSecondary{}).Where("chat_id = ? AND user_id <> ?", chat_id, user_id).First(&keys)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil, errors.New("ключ не найден")
	}

	return keys, nil
}

func CheckSecondaryKey(user_id uint, chat_id uint32) (*Model.KeysSecondary, error){
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


func GetPublicKey(chatId uint32) (*Model.Keys, error){

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