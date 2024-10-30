package main

import (
	"context"
	"fmt"
	chat_actions "mypackages/actions/chat"
	"mypackages/controllers"
	"mypackages/helpers"
	"mypackages/proto/chat"
	"strconv"
	"strings"

	// "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type DataStreamConnect struct{
	ChatId int
	UserID uint
	Stream chat.ChatGreeter_StreamGetMessagesServer
}

type chatServer struct {
	chat.UnimplementedChatGreeterServer
	Conns map[string]controllers.DataStreamConnect
}

// func ChatServer() *chatServer {
// 	return &chatServer{
// 		Conns: make(map[int]grpc.BidiStreamingServer[chat.StreamGetMessagesRequest, chat.StreamGetMessagesResponse]),
// 	}
// }

func (s *chatServer) CreateChat(ctx context.Context, in *chat.CreateRequestChat) (*chat.CreateResponseChat, error) {
	return controllers.CreateChat(ctx, in)
}

func (s *chatServer) StreamGetChat(in *chat.Empty, requestStream chat.ChatGreeter_StreamGetChatServer) error {
	return controllers.StreamGetChat(in,requestStream)
}

func (s *chatServer) GetUnSuccessChats(ctx context.Context, in *chat.Empty) (*chat.GetUnSuccessChatsResponse, error) {
	return controllers.GetUnSuccessChats(ctx, in)
}

func (s *chatServer) CreateSecondaryKey(ctx context.Context, in *chat.CreateSecondaryKeyRequest) (*chat.CreateSecondaryKeyResponse, error){
	return controllers.CreateSecondaryKey(ctx, in)
}

func (s *chatServer) GetSecondaryKey(ctx context.Context, in *chat.GetSecondaryKeyRequest) (*chat.GetSecondaryKeyResponse, error) {
	return controllers.GetSecondaryKey(ctx, in)
}

func (s *chatServer) GetPublicKey(ctx context.Context, in *chat.GetPublicKeyRequest) (*chat.GetPublicKeyResponse, error) {
	return controllers.GetPublicKey(ctx, in)
}

func (s *chatServer) AcceptChat(ctx context.Context, in *chat.AcceptChatRequest) (*chat.AcceptChatResponse, error) {
	return controllers.AcceptChat(ctx, in)
}

func (s *chatServer) DissalowChat(ctx context.Context, in *chat.DissalowChatRequest) (*chat.DissalowChatResponse, error) {
	return controllers.DissalowChat(ctx, in)
}

func (s *chatServer) GetMessages(ctx context.Context, in *chat.GetMessagesRequest) (*chat.GetMessagesResponse, error) {
	return controllers.GetMessages(ctx, in)
}

func (s *chatServer) StreamGetMessagesGeneral(in *chat.Empty, responseStream chat.ChatGreeter_StreamGetMessagesGeneralServer) error{
	return controllers.StreamGetMessagesGeneral(in, responseStream)
}

func (s *chatServer) StreamGetMessages(stream chat.ChatGreeter_StreamGetMessagesServer) error {
	fmt.Println(s.Conns)
	ctx := stream.Context()
	md, _ := metadata.FromIncomingContext(ctx)
	user, err := helpers.GetUserFormMd(ctx)
	channel := make(chan *chat.StreamGetMessagesResponse)
	if err!=nil {
		return status.Error(codes.PermissionDenied, "пользователь не найден")
	}
	chatId, err := strconv.Atoi( md["chat_id"][0]) 
	if err != nil {
		return status.Error(codes.PermissionDenied, "чат не найден ошибка")
	}
	err, chatUser := chat_actions.CheckChat(uint32(chatId), user.ID)
	
	if err != nil {
		return status.Error(codes.PermissionDenied, "чат не найден ошибка")
	}

	go func ()  {
		for{
			messageResponse := <-channel
			for _, v := range s.Conns {
				if v.ChatId == int(chatId) {
					v.Stream.Send(messageResponse)
				}
			}
		}
	}()
	jwtToken, _ := md["authorization"]
	jwtToken = strings.Split(jwtToken[0], " ")
	s.Conns[jwtToken[1]] = controllers.DataStreamConnect{
		ChatId:  chatUser.ChatID,
		UserID: user.ID,
		Stream: stream,
	}

	// defer controllers.CloseConnect(s.Conns, jwtToken[1])

	return controllers.StreamGetMessages(stream, s.Conns, chatId, int(user.ID), &channel, jwtToken[1])
}

