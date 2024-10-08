package main

import (
	"context"
	"mypackages/controllers"
	"mypackages/proto/chat"
)

type chatServer struct {
	chat.UnimplementedChatGreeterServer
}

func (s *chatServer) CreateChat(ctx context.Context, in *chat.CreateRequestChat) (*chat.CreateResponseChat, error) {
	return controllers.CreateChat(ctx, in)
}



func (s *chatServer) GetChat(ctx context.Context, in *chat.GetRequestChat) (*chat.GetResponseChat, error) {
	return controllers.GetChat(ctx, in)
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