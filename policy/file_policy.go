package policy

import (
	"github.com/R1kkass/GoCloudGRPC/db"
	Model "github.com/R1kkass/GoCloudGRPC/models"
)

func FolderPolicyID(folder_id uint32, user *Model.User) bool {

	if folder_id == 0 {
		return true
	}

	var folder Model.Folder

	r := db.DB.Model(&folder).First(&folder, "user_id=? AND id=?", user.ID, folder_id)

	if r.RowsAffected == 0 || r.Error != nil {
		return false
	}

	return true
}

func SpacePolicy(size uint32) bool {
	var user Model.User
	var file Model.File

	result := db.DB.Model(&file).Select("sum(size) as size").Where("user_id=?", user.ID).Group("user_id").Scan(&file)

	if result.Error != nil {
		return false
	}

	if file.Size+int(size) > 1024*1024*1024 {
		return false
	}
	return true
}

func CanGetMessages(chatId uint, userId uint32) bool {
	var chatUser *Model.ChatUser
	result := db.DB.Model(&Model.ChatUser{}).Where("chat_id = ? AND user_id = ?", chatId, userId).First(chatUser)

	if result.RowsAffected == 0 || result.Error != nil {
		return false
	}

	return true
}
