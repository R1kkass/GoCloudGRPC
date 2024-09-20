package main

import (
	"context"
	"mypackages/controllers"
	"mypackages/proto/auth"
)

type authServer struct{
	auth.UnimplementedAuthGreetServer
}

func (r *authServer) Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {
	return controllers.Login(ctx, in)
}

func (r *authServer) Registration(ctx context.Context, in *auth.RegistrationRequest) (*auth.RegistrationResponse, error) {
	return controllers.Registration(ctx, in)
}

func (r *authServer) DHConnect(ctx context.Context, in *auth.DHConnectRequest) (*auth.DHConnectResponse, error) {
	return controllers.DHConnect(ctx, in)
}

func (r *authServer) DHSecondConnect(ctx context.Context, in *auth.DHSecondConnectRequest) (*auth.DHSecondConnectResponse, error) {
	return controllers.DHSecondConnect(ctx, in)
}