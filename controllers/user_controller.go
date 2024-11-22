package controllers

import (
	"context"
	"strconv"

	"github.com/R1kkass/GoCloudGRPC/db"
	"github.com/R1kkass/GoCloudGRPC/helpers"
	Model "github.com/R1kkass/GoCloudGRPC/models"
	"github.com/R1kkass/GoCloudGRPC/proto/users"
)

type UsersServer struct {
	users.UnimplementedUsersGreetServer
}

func (s *UsersServer) GetUsers(ctx context.Context, in *users.GetUsersRequest) (*users.GetUsersResponse, error) {
	var usersList []Model.User
	user, _ := helpers.GetUserFormMd(ctx)

	db.DB.Model(&Model.User{}).Where("id != " + strconv.Itoa(int(user.ID)) + " AND (name LIKE name '%" + in.GetUserName() + "%' OR email LIKE '%" + in.GetUserName() + "%')").Find(&usersList)

	var usersResponse []*users.Users

	for i := 0; i < len(usersList); i++ {
		var u *users.Users = &users.Users{Id: int32(usersList[i].ID), Name: usersList[i].Name, Email: usersList[i].Email}
		usersResponse = append(usersResponse, u)
	}

	return &users.GetUsersResponse{Data: usersResponse}, nil
}
