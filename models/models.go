package Model

import (
	"time"

	"gorm.io/gorm"
)

type DefaultModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Folder struct {
	DefaultModel
	UserRelation
	FolderRelation

	NameFolder string    `json:"name_folder"`
	AccessId   int       `json:"access_id"`
	Access     *Accesses `json:"access" gorm:"default:null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type User struct {
	DefaultModel
	ID       uint   `gorm:"primary_key" json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
	Files    []File
	Folders  []Folder
}

type File struct {
	DefaultModel
	UserRelation
	FolderRelation

	Size         int       `json:"size"`
	FileName     string    `json:"file_name"`
	FileNameHash string    `json:"file_name_hash"`
	AccessId     int       `json:"access_id"`
	Access       *Accesses `json:"access" gorm:"default:null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Permissions struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	GuardName         string  `json:"guard_name"`
	PermissionGroupId int     `json:"permission_group_id"`
	Type              string  `json:"type"`
	CreatedAt         []uint8 `json:"created_at"`
	UpdatedAt         []uint8 `json:"updated_at"`
}

type PermissionGroups struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Title     string  `json:"title"`
	CreatedAt []uint8 `json:"created_at"`
	UpdatedAt []uint8 `json:"updated_at"`
}

type RoleHasPermissions struct {
	PermissionId int `json:"permission_id"`
	RoleId       int `json:"role_id"`
}

type ModelHasRoles struct {
	RoleId    int    `json:"role_id"`
	ModelType string `json:"model_type"`
	ModelId   int    `json:"model_id"`
}

type Accesses struct {
	DefaultModel
	Name string `json:"name"`
}

type RequestAccess struct {
	DefaultModel
	FileRelation
	FolderRelation

	UserID        int     `json:"user_id"`
	CurrentUserID int     `json:"current_user_id"`
	StatusID      int     `json:"status_id"`
	User          *User   `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CurrentUser   *User   `json:"current_user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status        *Status `json:"status" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Status struct {
	DefaultModel
	Name string `json:"name"`
}

type Chat struct {
	DefaultModel

	Messages []Message `json:"messages"`
	Message  Message   `json:"message" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ChatUsers []ChatUser `json:"chat_users"`
	NameChat  string     `json:"name_chat"`
}

type ChatUser struct {
	DefaultModel
	UserRelation

	ChatID int   `json:"chat_id"`
	Chat   *Chat `json:"chat" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SubmitCreate bool `json:"submit_create" gorm:"default:false"`
}

type Message struct {
	DefaultModel
	UserRelation

	ChatID   int    `json:"chat_id"`
	Chat     *Chat  `json:"chat" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Text     string `json:"text"`
}

type Keys struct {
	DefaultModel

	ChatID uint `json:"user_id"`
	P string `json:"p"`
	G int64 `json:"g"`
}

type KeysSecondary struct {
	DefaultModel

	UserID uint `json:"user_id"`
	ChatID uint `json:"chat_id"`
	Key string `json:"key"`
}

type SavedKeys struct {
	DefaultModel

	UserID uint `json:"user_id"`
	Token string `json:"token"`
	Ip string `json:"ip"`
	Name string `json:"name"`
	DateEnd uint `json:"date_end"`
}

type Notification struct {
	DefaultModel
	UserRelation

	Title string `json:"title"`
	Description string `json:"description"`
	Readed bool `json:"readed"`
}