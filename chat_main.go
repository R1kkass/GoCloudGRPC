package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mypackages/db"
	"mypackages/helpers"
	Model "mypackages/models"
	"mypackages/proto/chat"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type chatServer struct {
	chat.UnimplementedChatGreeterServer
}

func (s *chatServer) CreateChat(ctx context.Context, in *chat.CreateRequestChat) (*chat.CreateResponseChat, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	if err := FindUser(ctx, in); err != nil {
		return &chat.CreateResponseChat{
			Message: "Ошибка",
		}, status.Error(codes.NotFound, "Пользователь не найден")
	}

	var chats Model.Chat

	jwtToken, _ := md["authorization"]

	user, _ := helpers.GetUser(jwtToken)

	if err := CheckChat(ctx, in); err != nil {
		return &chat.CreateResponseChat{
			Message: "Ошибка",
		}, status.Error(codes.AlreadyExists, "Чат уже существует")
	}

	db.DB.Create(
		&chats,
	)
	var chatUsers Model.ChatUser

	db.DB.Model(&chatUsers).Create(
		&Model.ChatUser{
			ChatID: int(chats.ID),
			UserRelation: Model.UserRelation{
				UserID: int(user.ID),
			},
		},
	).Create(&Model.ChatUser{
		ChatID: int(chats.ID),
		UserRelation: Model.UserRelation{
			UserID: int(in.GetOtherId()),
		},
	})

	return &chat.CreateResponseChat{
		Message: fmt.Sprintf("%v", user),
	}, nil
}

func CheckChat(ctx context.Context, in *chat.CreateRequestChat) error {
	md, _ := metadata.FromIncomingContext(ctx)
	jwtToken, _ := md["authorization"]

	user, _ := helpers.GetUser(jwtToken)

	var usersChat Model.ChatUser

	result := db.DB.Raw(`SELECT count(*), chat_id from chat_users WHERE chat_id in (SELECT chat_id FROM chat_users Where user_id = ? INTERSECT SELECT chat_id FROM chat_users Where user_id = ?) GROUP BY chat_id`, user.ID, in.GetOtherId()).Scan(&usersChat)

	log.Println(result)
	if result.RowsAffected != 0 {
		return errors.New("чат уже существует")
	}

	return nil
}

func FindUser(ctx context.Context, in *chat.CreateRequestChat) error {
	var users Model.User
	result := db.DB.Model(&users).Where("id = ?", in.GetOtherId()).Find(&users)

	if result.RowsAffected == 0 {
		return errors.New("пользователь не найден")
	}

	return nil
}

func (s *chatServer) GetChat(ctx context.Context, in *chat.Empty) (*chat.GetResponseChat, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwtToken, _ := md["authorization"]

	user, _ := helpers.GetUser(jwtToken)

	var chats []*chat.ChatUsers

	db.DB.Model(&Model.ChatUser{}).Preload("User").Preload("Chat").Preload("Chat.Message").Preload("Chat.ChatUsers.User").Where("user_id = ?", user.ID).Find(&chats)
 
	out, _ := json.Marshal(chats)
	md.Set("text", string(out))
	
	return &chat.GetResponseChat{Chats: chats}, nil
}
