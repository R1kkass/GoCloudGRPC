package main

import (
	"context"

	"github.com/R1kkass/GoCloudGRPC/controllers"
	"github.com/R1kkass/GoCloudGRPC/proto/access"
)

type accessServer struct {
	access.UnimplementedAccessGreeterServer
}

func (s *accessServer) CreateAccess(ctx context.Context, in *access.RequestAccess) (*access.ResponseAccess, error) {
	return controllers.CreateAccess(ctx, in)
}

func (s *accessServer) GetAccesses(ctx context.Context, in *access.Empty) (*access.GetAccessesResponse, error) {
	return controllers.GetAccesses(ctx, in)
}

func (s *accessServer) ChangeAccess(ctx context.Context, in *access.ChangeAccessRequest) (*access.ChangeAccessResponse, error) {
	return controllers.ChangeAccess(ctx, in)
}
