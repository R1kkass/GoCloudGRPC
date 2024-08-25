package main

import (
	"context"
	"flag"
	"log"
	"net"

	"mypackages/consts"
	db "mypackages/db"
	"mypackages/interceptor"
	Model "mypackages/models"
	access "mypackages/proto/access"
	"mypackages/proto/chat"
	users "mypackages/proto/users"

	"github.com/joho/godotenv"

	"google.golang.org/grpc"
)

type Message struct {
	Email  string `json:"email"`
	Status string `json:"status"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

const (
	defaultName = "world"
)

func main() {
	db.ConnectDatabase()
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.CheckAuthInterceptor),
	)

	users.RegisterUsersGreetServer(s, &usersServer{})
	access.RegisterAccessGreeterServer(s, &accessServer{})
	chat.RegisterChatGreeterServer(s, &chatServer{})

	log.Printf("server listening at %v", lis.Addr())

	var chats []*chat.ChatUsers
	// var chats *Model.ChatUser

	// db.DB.Model(&Model.ChatUser{}).Preload("User").Preload("Chat.Message").Where("user_id = ?", user.ID).Find(&chats)
	// db.DB.Model(&Model.ChatUser{}).Preload("User").Preload("Chat").Preload("Chat.Message").Where("user_id = ? AND chat_id=1", 2).Find(&chats)
	db.DB.Model(&Model.ChatUser{}).Preload("User").Preload("Chat").Preload("Chat.Message").Preload("Chat.ChatUsers.User").Where("user_id = ?", 2).Find(&chats)

	// log.Printf("Some debug: %v", &chat.ChatUsers{
	// 	Chat: &chat.Chat{
	// 		Id:      uint32(chats.ID),
	// 		Message: &chat.Message{Id: uint32(chats.Chat.Message.ID), Text: chats.Chat.Message.Text},
	// 	},
	// })

	log.Printf("Some debug: %v", chats[0].Chat.ChatUsers[1].User)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

type usersServer struct {
	users.UnimplementedUsersGreetServer
}

func (s *usersServer) GetUsers(ctx context.Context, in *users.GetUsersRequest) (*users.GetUsersResponse, error) {
	var usersList []Model.User
	log.Printf("Received: %v", in.GetUserName())
	db.DB.Model(&Model.User{}).Where("name LIKE ?", "%"+in.GetUserName()+"%").Or("email LIKE ?", "%"+in.GetUserName()+"%").Find(&usersList)
	var usersResponse []*users.Users

	for i := 0; i < len(usersList); i++ {
		var u *users.Users = &users.Users{Id: int32(usersList[i].ID), Name: usersList[i].Name, Email: usersList[i].Email}
		usersResponse = append(usersResponse, u)
	}

	return &users.GetUsersResponse{Data: usersResponse}, nil
}

func (s *usersServer) GetContentUser(ctx context.Context, in *users.GetContentUserRequest) (*users.GetContentUserResponse, error) {
	var contentFiles []*users.File
	var contentFolders []*users.Folder
	var contentFoldersRequestAccess []*users.Folder
	var contentFileRequestAccess []*users.File

	db.DB.Model(&Model.File{}).Where("user_id = ? AND access_id = ?", in.GetId(), consts.OPEN).Find(&contentFiles)
	db.DB.Model(&Model.Folder{}).Where("user_id = ? AND access_id = ?", in.GetId(), consts.OPEN).Find(&contentFolders)
	db.DB.Model(&Model.Folder{}).Where("user_id = ? AND access_id = ?", in.GetId(), consts.WITH_PERMISSION).Find(&contentFoldersRequestAccess)
	db.DB.Model(&Model.File{}).Where("user_id = ? AND access_id = ?", in.GetId(), consts.WITH_PERMISSION).Find(&contentFileRequestAccess)

	return &users.GetContentUserResponse{
		Data: &users.Content{
			Files:               contentFiles,
			Folder:              contentFolders,
			FolderRequestAccess: contentFoldersRequestAccess,
			FileRequestAccess:   contentFileRequestAccess,
		},
	}, nil
}
