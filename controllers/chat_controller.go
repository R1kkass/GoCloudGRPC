package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	chat_actions "mypackages/actions/chat"
	"mypackages/db"
	"mypackages/helpers"
	Model "mypackages/models"
	"mypackages/proto/chat"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func CreateChat(ctx context.Context, in *chat.CreateRequestChat) (*chat.CreateResponseChat, error) {

	var chats Model.Chat

	user, err := helpers.GetUserFormMd(ctx)

	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Пользователь не найден")
	}

	if err := chat_actions.CheckChatExist(ctx, in); err != nil {
		return &chat.CreateResponseChat{
			Message: "Ошибка",
		}, status.Error(codes.AlreadyExists, "Чат уже существует")
	}

	var p string
	var g int64

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		result := db.DB.Create(
			&chats,
		)
		var chatUsers Model.ChatUser
		
		if result.Error != nil {
			return errors.New("чат не создан")
		}

		result = db.DB.Model(&chatUsers).Create(
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

		if result.Error != nil {
			return errors.New("чат не создан")
		}

		p, g = chat_actions.SendFirstParams(&chats);
		return nil
	})
	
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &chat.CreateResponseChat{
		Message: fmt.Sprintf("%v", user),
		Keys: &chat.Keys{
			P: p,
			G: g,
		},
		ChatId: uint32(chats.ID),
	}, nil
}



func GetChat(ctx context.Context, in *chat.Empty) (*chat.GetResponseChat, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwtToken, _ := md["authorization"]
	user, _ := helpers.GetUser(jwtToken)

	var chats []*chat.ChatUsers

	db.DB.Model(&Model.ChatUser{}).Preload("User").Preload("Chat").Preload("Chat.Message").Preload("Chat.ChatUsers.User").Where("user_id = ?", user.ID).Find(&chats)
 
	out, _ := json.Marshal(chats)
	md.Set("text", string(out))
	
	return &chat.GetResponseChat{Chats: chats}, nil
}


func CreateSecondaryKey(ctx context.Context, in *chat.CreateSecondaryKeyRequest) (*chat.CreateSecondaryKeyResponse, error){
	user, err := helpers.GetUserFormMd(ctx)

	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Пользователь не найден")
	}
	fmt.Println(user.ID)
	_, err = chat_actions.CheckSecondaryKey(user.ID, in.GetChatId())

	if err != nil {
		return nil, status.Error(codes.AlreadyExists, "Ключ уже создан")
	}

	db.DB.Create(&Model.KeysSecondary{
		UserID: user.ID,
		ChatID: uint(in.GetChatId()),
		Key: in.GetKey(),
	})
	
	return &chat.CreateSecondaryKeyResponse{
		Message: "Успешно",
	}, nil
}

func GetSecondaryKey(ctx context.Context, in *chat.GetSecondaryKeyRequest) (*chat.GetSecondaryKeyResponse, error) {
	
	user, err := helpers.GetUserFormMd(ctx)
	fmt.Println(user.ID)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Пользователь не найден")
	}

	err, _ = chat_actions.CheckChat(in.GetChatId(), user.ID)

	if err != nil {
		return nil, status.Error(codes.NotFound, "Чат не найден")
	}

	key, err := chat_actions.GetSecondaryKey(user.ID, in.GetChatId())
	
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	
	keys, err := chat_actions.GetPublicKey(in.GetChatId())

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	
	return &chat.GetSecondaryKeyResponse{
		Key: key.Key,
		P: keys.P,
	}, nil
}

func GetPublicKey(ctx context.Context, in *chat.GetPublicKeyRequest) (*chat.GetPublicKeyResponse, error) {

	if err := helpers.CheckChat(ctx, in.GetChatId()); err != nil {
		return nil, status.Error(codes.NotFound, "Чат не найден")
	}

	keys, err := chat_actions.GetPublicKey(in.GetChatId())
	
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &chat.GetPublicKeyResponse{
		G: keys.G,
		P: keys.P,
	}, nil
}
