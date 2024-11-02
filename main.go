package main

import (
	"log"
	"net"

	db "github.com/R1kkass/GoCloudGRPC/db"
	"github.com/R1kkass/GoCloudGRPC/interceptor"
	access "github.com/R1kkass/GoCloudGRPC/proto/access"
	"github.com/R1kkass/GoCloudGRPC/proto/auth"
	"github.com/R1kkass/GoCloudGRPC/proto/chat"
	"github.com/R1kkass/GoCloudGRPC/proto/files"
	"github.com/R1kkass/GoCloudGRPC/proto/keys"
	"github.com/R1kkass/GoCloudGRPC/proto/notification"
	users "github.com/R1kkass/GoCloudGRPC/proto/users"
	"github.com/R1kkass/GoCloudGRPC/structs"

	"github.com/joho/godotenv"

	"google.golang.org/grpc"
	// "github.com/R1kkass/GoCloudGRPC/tls"
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
		// grpc.StreamInterceptor(interceptor.CheckAuthInterceptorStream),
	)

	users.RegisterUsersGreetServer(s, &usersServer{})
	access.RegisterAccessGreeterServer(s, &accessServer{})
	chat.RegisterChatGreeterServer(s, &chatServer{
		Conns: make(map[string]structs.DataStreamConnect),
	})
	auth.RegisterAuthGreetServer(s, &authServer{})
	keys.RegisterKeysGreeterServer(s, &keysServer{})
	files.RegisterFilesGreeterServer(s, &filesServer{})
	notification.RegisterNotificationGreeterServer(s, &notificationServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
