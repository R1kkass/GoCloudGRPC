package interceptor

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var seccretKey, _ = os.LookupEnv("SECRET_KEY")

func CheckAuth(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)

	jwtToken, ok := md["authorization"]
	log.Printf("Received: %v", md)
	if !ok || len(jwtToken) < 1 {
		return status.Error(codes.Unauthenticated, "не авторизован")
	}

	jwtToken = strings.Split(jwtToken[0], " ")

	token, err := jwt.Parse(jwtToken[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(seccretKey), nil
	})

	if !token.Valid || err != nil {
		return status.Error(codes.Unauthenticated, "не авторизован")
	}
	return nil
}

func CheckAuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	jwtToken, ok := md["authorization"]
	if !ok || len(jwtToken) < 1 {
		return nil, status.Error(codes.Unauthenticated, "не авторизован")
	}

	jwtToken = strings.Split(jwtToken[0], " ")

	token, err := jwt.Parse(jwtToken[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(seccretKey), nil
	})

	if !token.Valid || err != nil {
		return nil, status.Error(codes.Unauthenticated, "не авторизован")
	}
	m, err := handler(ctx, req)
	return m, err
}
