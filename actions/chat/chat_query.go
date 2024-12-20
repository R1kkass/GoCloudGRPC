package chat_actions

import (
	"errors"

	"github.com/R1kkass/GoCloudGRPC/db"
	Model "github.com/R1kkass/GoCloudGRPC/models"
	"github.com/R1kkass/GoCloudGRPC/proto/chat"
	"gorm.io/gorm"
)

func CreateChatTransaction(userId uint, companionId uint) (string, int64, *Model.Chat, error) {
	var p string
	var g int64
	var chat *Model.Chat
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(
			&chat,
		)
		var chatUsers Model.ChatUser

		if result.Error != nil {
			return errors.New("чат не создан")
		}

		result = tx.Model(&chatUsers).Create(
			&Model.ChatUser{
				ChatRelations: Model.ChatRelations{
					ChatID: chat.ID,
				},
				UserRelation: Model.UserRelation{
					UserID: userId,
				},
				SubmitCreate: Model.WaitChat,
			},
		).Create(&Model.ChatUser{
			ChatRelations: Model.ChatRelations{
				ChatID: chat.ID,
			},
			UserRelation: Model.UserRelation{
				UserID: companionId,
			},
			SubmitCreate: Model.UnSuccessChat,
		})
		if result.Error != nil {
			return errors.New("чат не создан")
		}

		p, g = SendFirstParams(chat)
		return nil
	})

	return p, g, chat, err
}

func UploadChatFileTransaction(userId uint, chatId uint, messageId uint, fileName string) (*Model.ChatFile, error) {
	var chatFile *Model.ChatFile
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		chatFile = &Model.ChatFile{
			ChatRelations: Model.ChatRelations{
				ChatID: chatId,
			},
			MessageRelations: Model.MessageRelations{
				MessageID: messageId,
			},
			UserRelation: Model.UserRelation{
				UserID: userId,
			},
			FileName: fileName,
		}
		r := tx.Create(&chatFile)
		if r.RowsAffected == 0 || r.Error != nil {
			return errors.New("ошибка создания ChatFile")
		}
		return nil
	})

	return chatFile, err
}

func DissalowChatTransaction(chatId uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("chat_id = ?", chatId).Unscoped().Delete(&Model.ChatUser{})

		if result.Error != nil {
			return errors.New("ошибка")
		}

		result = tx.Where("id=?", chatId).Unscoped().Delete(&Model.Chat{})
		if result.Error != nil {
			return errors.New("ошибка")
		}

		return nil
	})
}

func StreamGetChat(userId uint) ([]*chat.ChatUsersCount, error) {
	var chats []*chat.ChatUsersCount
	r := db.DB.Model(&Model.ChatUser{}).Select(`chat_users.*, COALESCE(messages.created_at,'2022-10-19 15:23:53.252567+00') as create_at_message, count(un_readed_messages.id) as un_readed_messages_count`).
		Joins("LEFT JOIN un_readed_messages ON un_readed_messages.chat_id = chat_users.chat_id AND un_readed_messages.user_id = ?", userId).
		Joins("LEFT JOIN (SELECT * FROM (SELECT distinct on(chat_id) chat_id, created_at, status_message FROM messages WHERE status_message = 'success' ORDER BY chat_id, created_at DESC) t ORDER BY created_at DESC) AS messages ON messages.chat_id = chat_users.chat_id AND messages.status_message = 'success'", userId).
		Preload("User").Preload("Chat").Preload("Chat.ChatUsers.User").
		Preload("Chat.Message", func(db *gorm.DB) *gorm.DB {
			return db.Where("status_message = 'success'").Order("messages.id ASC")
		}).
		Preload("Chat.Message.User").
		Preload("Chat.Message.ChatFiles").
		Where("chat_users.user_id = ? AND submit_create = ?", userId, Model.CreatedChat).
		Group("chat_users.id, un_readed_messages.chat_id, messages.created_at").
		Order("create_at_message DESC, chat_users.created_at DESC").
		Find(&chats)

	if r.Error != nil {
		return nil, r.Error
	}

	return chats, nil
}
