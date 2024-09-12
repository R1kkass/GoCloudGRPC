package main

import (
	"context"
	"mypackages/controllers"
	"mypackages/proto/auth"
)

type authServer struct{
	auth.UnimplementedAuthGreetServer
}

func (s *authServer) Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {
	return controllers.Login(ctx, in)
}

func (s *authServer) Registration(ctx context.Context, in *auth.RegistrationRequest) (*auth.RegistrationResponse, error) {
	return controllers.Registration(ctx, in)
}

func (s *authServer) DHConnect(ctx context.Context, in *auth.DHConnectRequest) (*auth.DHConnectResponse, error) {
	return controllers.DHConnect(ctx, in)
}

func (s *authServer) DHSecondConnect(ctx context.Context, in *auth.DHSecondConnectRequest) (*auth.DHSecondConnectResponse, error) {
	return controllers.DHSecondConnect(ctx, in)
}