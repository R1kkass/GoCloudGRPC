package main

import (
	"context"
	"mypackages/controllers"
	users "mypackages/proto/users"
)


type usersServer struct {
	users.UnimplementedUsersGreetServer
}

func (s *usersServer) GetUsers(ctx context.Context, in *users.GetUsersRequest) (*users.GetUsersResponse, error) {
	return controllers.GetUsers(ctx, in)
}

func (s *usersServer) GetContentUser(ctx context.Context, in *users.GetContentUserRequest) (*users.GetContentUserResponse, error) {
	return controllers.GetContentUser(ctx, in)
}