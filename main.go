package main

import (
	"log"
	"net"

	db "mypackages/db"
	"mypackages/interceptor"
	access "mypackages/proto/access"
	"mypackages/proto/auth"
	"mypackages/proto/chat"
	"mypackages/proto/files"
	"mypackages/proto/keys"
	"mypackages/proto/notification"
	users "mypackages/proto/users"

	"github.com/joho/godotenv"

	"google.golang.org/grpc"
	// "mypackages/tls"
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

func main() {
	db.ConnectDatabase()
	db.ConnectRedis()
	db.ConnectRedisNotification()

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// tlsCreds, err := tls.GenerateTLSCreds()
	if err != nil {
	   log.Fatal(err)
	}

	s := grpc.NewServer(
		// grpc.Creds(tlsCreds),
		grpc.UnaryInterceptor(interceptor.CheckAuthInterceptor),
	)

	users.RegisterUsersGreetServer(s, &usersServer{})
	access.RegisterAccessGreeterServer(s, &accessServer{})
	chat.RegisterChatGreeterServer(s, &chatServer{})
	auth.RegisterAuthGreetServer(s, &authServer{})
	keys.RegisterKeysGreeterServer(s, &keysServer{})
	files.RegisterFilesGreeterServer(s, &filesServer{})
	notification.RegisterNotificationGreeterServer(s, &notificationServer{})
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}