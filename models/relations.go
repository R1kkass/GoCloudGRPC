package Model

type UserRelation struct {
	UserID int   `json:"user_id"`
	User   *User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type FileRelation struct {
	File   *File `json:"file" gorm:"default:null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	FileID int   `json:"file_id" gorm:"default:null;"`
}

type FolderRelation struct {
	FolderID int     `json:"folder_id" gorm:"default:null;"`
	Folder   *Folder `json:"folder" gorm:"default:null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
