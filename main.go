package main

import (
	"log"
	"net"

	"github.com/R1kkass/GoCloudGRPC/controllers"
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
	db.AwsConnect()
	db.Migration()

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
		grpc.StreamInterceptor(interceptor.CheckAuthInterceptorStream),
		grpc.MaxSendMsgSize(1e+7),
		grpc.MaxRecvMsgSize(1e+7),
	)

	users.RegisterUsersGreetServer(s, &controllers.UsersServer{})
	access.RegisterAccessGreeterServer(s, &controllers.AccessServer{})
	chat.RegisterChatGreeterServer(s, &controllers.ChatServer{
		Conns: make(map[string]structs.DataStreamConnect),
	})
	auth.RegisterAuthGreetServer(s, &controllers.AuthServer{})
	keys.RegisterKeysGreeterServer(s, &controllers.KeysServer{})
	files.RegisterFilesGreeterServer(s, &controllers.FilesServer{})
	notification.RegisterNotificationGreeterServer(s, &controllers.NotificationServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
