package interceptor

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)


var seccretKey, _ = os.LookupEnv("SECRET_KEY")


func CheckAuthInterceptorStream(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, _ := metadata.FromIncomingContext(ss.Context())
	jwtToken, ok := md["authorization"]
	if !ok || len(jwtToken) < 1 {
		return status.Error(codes.Unauthenticated, "не авторизован")
	}

	jwtToken = strings.Split(jwtToken[0], " ")

	token, err := jwt.Parse(jwtToken[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(seccretKey), nil
	})

	if err != nil || !token.Valid {
		return status.Error(codes.Unauthenticated, "не авторизован")
	}
	err = handler(srv, ss)
	if err != nil {
		fmt.Println("RPC failed with error: %v", err)
	}

	return err
}

func CheckAuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if strings.Split(info.FullMethod, "/")[1]!="auth.AuthGreet" || info.FullMethod == "/auth.AuthGreet/CheckAuth" {


		md, _ := metadata.FromIncomingContext(ctx)
		jwtToken, ok := md["authorization"]
		if !ok || len(jwtToken) < 1 {
			return nil, status.Error(codes.Unauthenticated, "не авторизован")
		}
	
		jwtToken = strings.Split(jwtToken[0], " ")
	
		token, err := jwt.Parse(jwtToken[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(seccretKey), nil
		})
	
		if err != nil || !token.Valid {
			return nil, status.Error(codes.Unauthenticated, "не авторизован")
		}
	}
	m, err := handler(ctx, req)

	return m, err
}

