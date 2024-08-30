package helpers

import (
	"context"
	"errors"
	"mypackages/db"
	Model "mypackages/models"
)

func CheckChat(ctx context.Context, chatId uint32) error {
	user, err := GetUserFormMd(ctx)

	if err!=nil {
		return err
	}

	var usersChat Model.ChatUser

	result := db.DB.Model(&Model.ChatUser{}).Where("user_id = ? AND chat_id = ?", user.ID, chatId).First(&usersChat)

	if result.RowsAffected == 0 {
		return errors.New("пользователь не найден")
	}

	return nil
}