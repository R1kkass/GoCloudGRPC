package main

import (
	"log"
	"net"

	db "mypackages/db"
	"mypackages/interceptor"
	access "mypackages/proto/access"
	"mypackages/proto/auth"
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
	auth.RegisterAuthGreetServer(s, &authServer{})
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
