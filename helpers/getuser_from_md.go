package helpers

import (
	"context"
	"errors"
	Model "mypackages/models"

	"google.golang.org/grpc/metadata"
)

func GetUserFormMd(ctx context.Context) (*Model.User, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwtToken, _ := md["authorization"]
	user, err := GetUser(jwtToken)

	if err!=nil {
		return nil,  errors.New("пользователь не найден")
	}

	return user, nil
}