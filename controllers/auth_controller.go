package controllers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	auth_action "mypackages/actions/auth"
	"mypackages/db"
	"mypackages/helpers"
	Model "mypackages/models"
	"mypackages/proto/auth"
	"os"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type KeyUser struct {
	p *big.Int
	g int64
	b *big.Int
}

var keysUser  = make(map[string]KeyUser)
var generatedKeys  = make(map[string]string)

func Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {


	var user *Model.User;
	per, _ := peer.FromContext(ctx)
	ip := per.Addr.String()

	defer delete(generatedKeys, ip);

	secretKey := generatedKeys[ip]
	password := helpers.Decrypt(in.GetPassword(), secretKey)
	email := helpers.Decrypt(in.GetEmail(), secretKey)

	r := db.DB.Unscoped().Model(&user).First(&user, "email = ?", email)
	match := auth_action.CheckPasswordHash(password, user.Password)
	
	if r.RowsAffected == 0 || !match{
		return nil, status.Error(codes.Unknown, "Не авторизован")
	}


	secretToken, err := helpers.GenerateJWTToken(user, secretKey)

	if err != nil {
		return nil, status.Error(codes.Aborted, "Токен не создан")
	}

	db.DB.Model(&Model.SavedKeys{}).Create(&Model.SavedKeys{
		UserID: user.ID,
		Name: "",
		Ip: ip,
		Token: *secretToken,
	})

	return &auth.LoginResponse{
		AccessToken: *secretToken,
	}, nil
}

func Registration(ctx context.Context, in *auth.RegistrationRequest) (*auth.RegistrationResponse, error) {
	var user *Model.User;
	p, _ := peer.FromContext(ctx)
	ip:=p.Addr.String()
	
	secretKey := generatedKeys[ip]

	defer delete(generatedKeys, ip);

	password := helpers.Decrypt(in.GetPassword(), secretKey)
	email := helpers.Decrypt(in.GetEmail(), secretKey)
	name := helpers.Decrypt(in.GetName(), secretKey)

	r := db.DB.Unscoped().Model(&user).First(&user, "email = ?", email)
	if r.RowsAffected > 0 {
		return nil, status.Error(codes.AlreadyExists, "Такой пользователь уже есть")
	}

	pass, _ := auth_action.HashPassword(password)

	newUser := Model.User{
		Email: email,
		Password: pass,	
		Name: name,
	}

	r = db.DB.Create(&newUser)
	os.Mkdir("files/"+strconv.Itoa(int(user.ID)), os.ModePerm)

	if r.RowsAffected==0{
		return nil, status.Error(codes.Unauthenticated, "Не зарегистрирован")
	}

	secretToken, err := helpers.GenerateJWTToken(user, secretKey)

	if err != nil {
		return nil, status.Error(codes.Aborted, "Токен не создан")
	}

	db.DB.Model(&Model.SavedKeys{}).Create(&Model.SavedKeys{
		UserID: user.ID,
		Name: "",
		Ip: ip,
		Token: *secretToken,
	})

	return &auth.RegistrationResponse{
		AccessToken: *secretToken,
	}, nil
}

func DHConnect(ctx context.Context, in *auth.DHConnectRequest) (*auth.DHConnectResponse, error) {
	per, _ := peer.FromContext(ctx)
	ip := per.Addr.String()
	
	p,g := helpers.GeneratePubKeys()
	B, b, _ := helpers.GeneratePubKey(p, int(g))
	keysUser[ip] = KeyUser{
		p: p,
		g: g,
		b: b,
	}
	
	return &auth.DHConnectResponse{P: p.String(), G: g, B: B.String()}, nil
}

func DHSecondConnect(ctx context.Context, in *auth.DHSecondConnectRequest) (*auth.DHSecondConnectResponse, error) {
	p, _ := peer.FromContext(ctx)
	ip:=p.Addr.String()

	keys := keysUser[ip]
	b, _ := helpers.GenerateSecretKey(keys.p, keys.b.String(), in.GetA())
	
	bytes := []byte(b.String())
	hash := sha256.Sum256(bytes)
	hashString := hex.EncodeToString(hash[:])

	generatedKeys[ip] = string(hashString)[0:32]
	fmt.Println(string(hashString)[0:32])
	fmt.Println("string(hashString)[0:32]")
	return &auth.DHSecondConnectResponse{
		Message: "Ключ успешно создан",
	}, nil
}
