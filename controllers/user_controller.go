package controllers

import (
	"context"
	"mypackages/consts"
	"mypackages/db"
	"mypackages/helpers"
	Model "mypackages/models"
	"mypackages/proto/users"
	"strconv"
)

func GetUsers(ctx context.Context, in *users.GetUsersRequest) (*users.GetUsersResponse, error) {
	var usersList []Model.User
	user, _ := helpers.GetUserFormMd(ctx)

	db.DB.Model(&Model.User{}).Where("id != " + strconv.Itoa(int(user.ID)) +" AND (name LIKE name '%"+in.GetUserName()+"%' OR email LIKE '%"+in.GetUserName()+"%')").Find(&usersList)
	
	var usersResponse []*users.Users

	for i := 0; i < len(usersList); i++ {
		var u *users.Users = &users.Users{Id: int32(usersList[i].ID), Name: usersList[i].Name, Email: usersList[i].Email}
		usersResponse = append(usersResponse, u)
	}

	return &users.GetUsersResponse{Data: usersResponse}, nil
}

func GetContentUser(ctx context.Context, in *users.GetContentUserRequest) (*users.GetContentUserResponse, error) {

	var contentFilesChan = make(chan []*users.File, 1)
	var contentFoldersChan = make(chan []*users.Folder, 1)
	var contentFoldersRequestAccessChan = make(chan []*users.Folder, 1)
	var contentFileRequestAccessChan = make(chan []*users.File, 1)

	go func (){
		var contentFiles []*users.File
		db.DB.Model(&Model.File{}).Where("user_id = ? AND access_id = ?", in.GetId(), consts.OPEN).Find(&contentFiles)
		contentFilesChan <- contentFiles
	}()
	go func() {
		var contentFolders []*users.Folder
		db.DB.Model(&Model.Folder{}).Where("user_id = ? AND access_id = ?", in.GetId(), consts.OPEN).Find(&contentFolders)
		contentFoldersChan <- contentFolders
	}()
	go func(){
		var contentFoldersRequestAccess []*users.Folder
		db.DB.Model(&Model.Folder{}).Where("user_id = ? AND access_id = ?", in.GetId(), consts.WITH_PERMISSION).Find(&contentFoldersRequestAccess)
		contentFoldersRequestAccessChan <- contentFoldersRequestAccess
	}()
	go func() {
		var contentFileRequestAccess []*users.File
		db.DB.Model(&Model.File{}).Where("user_id = ? AND access_id = ?", in.GetId(), consts.WITH_PERMISSION).Find(&contentFileRequestAccess)
		contentFileRequestAccessChan <- contentFileRequestAccess
	}()

	return &users.GetContentUserResponse{
		Data: &users.Content{
			Files:               <- contentFilesChan,
			Folder:              <- contentFoldersChan,
			FolderRequestAccess: <- contentFoldersRequestAccessChan,
			FileRequestAccess:   <- contentFileRequestAccessChan,
		},
	}, nil
}