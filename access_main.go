package main

import (
	"context"
	"mypackages/consts"
	"mypackages/db"
	"mypackages/helpers"
	Model "mypackages/models"
	"mypackages/proto/access"
	"os"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var seccretKey, _ = os.LookupEnv("SECRET_KEY")

func (s *accessServer) CreateAccess(ctx context.Context, in *access.RequestAccess) (*access.ResponseAccess, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	var file Model.File;
	var folder Model.Folder;
	jwtToken, _ := md["authorization"]

	jwtToken = strings.Split(jwtToken[0], " ")
	user, ok := helpers.GetUser(jwtToken[1])

	if !ok {
		return nil, status.Error(codes.NotFound, "Пользователь не найден");
	}

	db.DB.Model(&Model.Folder{}).Where("id = ? AND user_id = ?", in.GetFolderId(), in.GetUserId()).First(&folder)
	db.DB.Model(&Model.File{}).Where("id = ? AND user_id = ?", in.GetFileId(), in.GetUserId()).First(&file)

	var request_access Model.RequestAccess

	result := db.DB.Model(&Model.RequestAccess{}).Where("file_id = ? AND folder_id = ? AND user_id = ? AND current_user_id = ?", in.GetFileId(), in.GetFolderId(), in.GetUserId(), user.ID, ).Find(&request_access)

	if in.GetUserId() != int32(user.ID) && 
		(folder.AccessId==consts.WITH_PERMISSION || folder.AccessId==consts.WITH_PERMISSION) && result.RowsAffected==0{
			db.DB.Create(&Model.RequestAccess{
				UserID: int(in.GetUserId()),
				CurrentUserID: int(user.ID),
				FileID: int(in.GetFileId()),
				FolderID: int(in.GetFolderId()),
				StatusID: consts.EXPECTATION,
			});
	} else {
		return nil, status.Error(codes.ResourceExhausted, "Пользователь не найден");
	}

	return &access.ResponseAccess{
		Message: "success",
	}, nil
}

func (s *accessServer) GetAccesses(ctx context.Context, in *access.Empty) (*access.GetAccessesResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	var request_access []*access.RequestAccessData;
	jwtToken, _ := md["authorization"]
	jwtToken = strings.Split(jwtToken[0], " ")
	user, ok := helpers.GetUser(jwtToken[1])

	if !ok {
		return nil, status.Error(codes.NotFound, "Пользователь не найден");
	}

	db.DB.Model(&Model.RequestAccess{}).Where("user_id = ?", user.ID).Find(&request_access)
	
	return &access.GetAccessesResponse{
		Accesses: request_access,
	}, nil
}

func(s *accessServer) ChangeAccess(ctx context.Context, in *access.ChangeAccessRequest) (*access.ChangeAccessResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	
	jwtToken, _ := md["authorization"]
	jwtToken = strings.Split(jwtToken[0], " ")
	user, ok := helpers.GetUser(jwtToken[1])

	if !ok {
		return nil, status.Error(codes.NotFound, "Пользователь не найден");
	}

	result := db.DB.Model(&Model.RequestAccess{}).Where("id = ? AND user_id = ?", in.GetId(), user.ID).Update("status_id = ", in.GetStatus())

	if result.RowsAffected==0 {
		return nil, status.Error(codes.NotFound, "Файл не найден")
	}

	return &access.ChangeAccessResponse{
		Message: "success",
	}, nil
}