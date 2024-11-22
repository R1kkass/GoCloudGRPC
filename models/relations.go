package Model

type UserRelation struct {
	UserID uint   `json:"user_id"`
	User   *User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type FileRelation struct {
	File   *File `json:"file" gorm:"default:null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	FileID uint   `json:"file_id" gorm:"default:null;"`
}

type FolderRelation struct {
	FolderID uint     `json:"folder_id" gorm:"default:null;"`
	Folder   *Folder `json:"folder" gorm:"default:null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ChatRelations struct {
	ChatID uint   `json:"chat_id"`
	Chat   *Chat `json:"chat" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type MessageRelations struct {
	MessageID uint     `json:"message_id"`
	Message   *Message `json:"message" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
