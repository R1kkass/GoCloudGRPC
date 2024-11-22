package controllers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"

	auth_action "github.com/R1kkass/GoCloudGRPC/actions/auth"
	"github.com/R1kkass/GoCloudGRPC/db"
	"github.com/R1kkass/GoCloudGRPC/helpers"
	Model "github.com/R1kkass/GoCloudGRPC/models"
	"github.com/R1kkass/GoCloudGRPC/proto/auth"
	"github.com/R1kkass/GoCloudGRPC/validate"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type KeyUser struct {
	p *big.Int
	g int64
	b *big.Int
}

type AuthServer struct {
	auth.UnimplementedAuthGreetServer
}

var keysUser = make(map[string]KeyUser)

func (s *AuthServer) Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error Login: ", r)
		}
	}()

	err := validate.Valid(
		validate.ValidType{
			"email": validate.ValidateStruct{
				Rule:  "required",
				Value: in.GetEmail(),
			},
			"password": validate.ValidateStruct{
				Rule:  "required",
				Value: in.GetPassword(),
			},
		},
	)

	if err != nil{
		return nil, status.Error(codes.Unknown, "Неизвестная ошибка")
	}

	var user *Model.User
	per, _ := peer.FromContext(ctx)
	ip := per.Addr.String()

	secretKey, err := db.ConnectRedisDB.HGet(ctx, "generatedKeys:", ip).Result()

	if err != nil {
		return nil, status.Error(codes.Aborted, "Ключи не созданы")
	}

	password := helpers.Decrypt(in.GetPassword(), secretKey)
	email := helpers.Decrypt(in.GetEmail(), secretKey)

	r := db.DB.Unscoped().Model(&user).First(&user, "email = ?", email)
	match := auth_action.CheckPasswordHash(password, user.Password)

	if r.RowsAffected == 0 || !match {
		return nil, status.Error(codes.Unknown, "Не авторизован")
	}

	secretToken, err := helpers.GenerateJWTToken(user, secretKey)

	if err != nil {
		return nil, status.Error(codes.Aborted, "Токен не создан")
	}

	db.DB.Model(&Model.SavedKeys{}).Create(&Model.SavedKeys{
		UserID: user.ID,
		Name:   "",
		Ip:     ip,
		Token:  *secretToken,
	})

	return &auth.LoginResponse{
		AccessToken: *secretToken,
	}, nil
}

func (s *AuthServer) Registration(ctx context.Context, in *auth.RegistrationRequest) (*auth.RegistrationResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error Registration: ", r)
		}
	}()

	err := validate.Valid(
		validate.ValidType{
			"email": validate.ValidateStruct{
				Rule:  "required",
				Value: in.GetEmail(),
			},
			"password": validate.ValidateStruct{
				Rule:  "required",
				Value: in.GetPassword(),
			},
			"name": validate.ValidateStruct{
				Rule:  "required|min:2",
				Value: in.GetName(),
			},
		},
	)

	if err != nil{
		return nil, status.Error(codes.Unknown, "Неизвестная ошибка")
	}

	var user *Model.User
	p, _ := peer.FromContext(ctx)
	ip := p.Addr.String()

	secretKey, err := db.ConnectRedisDB.HGet(ctx, "generatedKeys:", ip).Result()

	if err != nil {
		return nil, status.Error(codes.Aborted, "Ключ не созданы")
	}

	password := helpers.Decrypt(in.GetPassword(), secretKey)
	email := helpers.Decrypt(in.GetEmail(), secretKey)
	name := helpers.Decrypt(in.GetName(), secretKey)

	r := db.DB.Unscoped().Model(&user).First(&user, "email = ?", email)
	if r.RowsAffected > 0 {
		return nil, status.Error(codes.AlreadyExists, "Такой пользователь уже есть")
	}

	pass, _ := auth_action.HashPassword(password)

	newUser := Model.User{
		Email:    email,
		Password: pass,
		Name:     name,
	}

	r = db.DB.Create(&newUser).First(&newUser)
	var pathFileFolder, _ = os.LookupEnv("PATH_FILES")

	os.Mkdir(pathFileFolder+strconv.Itoa(int(newUser.ID)), os.ModePerm)
	var pathKeyFolder, _ = os.LookupEnv("PATH_KEYS")

	os.Create(pathKeyFolder + strconv.Itoa(int(newUser.ID)))
	if r.RowsAffected == 0 {
		return nil, status.Error(codes.Unauthenticated, "Не зарегистрирован")
	}

	secretToken, err := helpers.GenerateJWTToken(&newUser, secretKey)

	if err != nil {
		return nil, status.Error(codes.Aborted, "Токен не создан")
	}

	db.DB.Model(&Model.SavedKeys{}).Create(&Model.SavedKeys{
		UserID: newUser.ID,
		Name:   "",
		Ip:     ip,
		Token:  *secretToken,
	})

	return &auth.RegistrationResponse{
		AccessToken: *secretToken,
	}, nil
}

func (s *AuthServer) DHConnect(ctx context.Context, in *auth.DHConnectRequest) (*auth.DHConnectResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error DHConnect: ", r)
		}
	}()
	
	per, _ := peer.FromContext(ctx)
	ip := per.Addr.String()

	p, g := helpers.GeneratePubKeys()
	B, b, _ := helpers.GeneratePubKey(p, int(g))
	keysUser[ip] = KeyUser{
		p: p,
		g: g,
		b: b,
	}

	keys := map[string]any{
		"p": p.String(),
		"b": b.String(),
	}

	err := db.ConnectRedisDB.HMSet(ctx, "keysUser:"+ip, keys).Err()

	if err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Aborted, "Ключ не создан")
	}

	return &auth.DHConnectResponse{P: p.String(), G: g, B: B.String()}, nil
}

func (s *AuthServer) DHSecondConnect(ctx context.Context, in *auth.DHSecondConnectRequest) (*auth.DHSecondConnectResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error DHSecondConnect: ", r)
		}
	}()
	
	p, _ := peer.FromContext(ctx)
	ip := p.Addr.String()

	// keys := keysUser[ip]
	KeyP, err := db.ConnectRedisDB.HMGet(ctx, "keysUser:"+ip, "p").Result()

	if err != nil {
		return nil, status.Error(codes.Aborted, "Ключ не создан")
	}

	KeyB, err := db.ConnectRedisDB.HMGet(ctx, "keysUser:"+ip, "b").Result()
	if err != nil {
		return nil, status.Error(codes.Aborted, "Ключ не создан")
	}

	b, _ := helpers.GenerateSecretKey(fmt.Sprintf("%v", KeyP[0]), fmt.Sprintf("%v", KeyB[0]), in.GetA())

	bytes := []byte(b.String())
	hash := sha256.Sum256(bytes)
	hashString := hex.EncodeToString(hash[:])
	err = db.ConnectRedisDB.HSet(ctx, "generatedKeys:", ip, string(hashString)[0:32]).Err()

	if err != nil {
		return nil, status.Error(codes.Aborted, "Ключ не создан")
	}

	return &auth.DHSecondConnectResponse{
		Message: "Ключ успешно создан",
	}, nil
}

func (s *AuthServer) CheckAuth(ctx context.Context, in *auth.Empty) (*auth.Empty, error) {
	return &auth.Empty{}, nil
}