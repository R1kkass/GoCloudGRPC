package controllers

import (
	"context"
	"math/big"
	auth_action "mypackages/actions/auth"
	"mypackages/db"
	"mypackages/helpers"
	Model "mypackages/models"
	"mypackages/proto/auth"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var secretKey, _ = os.LookupEnv("SECRET_KEY")
var jwtSecretKey = []byte(secretKey)

type KeyUser struct {
	p *big.Int
	g int64
	b *big.Int
}

var keysUser  = make(map[string]KeyUser)
var generatedKeys  = make(map[string]string)

func Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {
	var user *Model.User;

	r := db.DB.Unscoped().Model(&user).First(&user, "email = ?", in.GetEmail())
	match := auth_action.CheckPasswordHash(in.GetPassword(), user.Password)
	
	if r.RowsAffected == 0 || !match{
		return nil, status.Error(codes.Unknown, "Не авторизован")
	}

    payload := jwt.MapClaims{
		"email": user.Email,
        "exp": time.Now().Add(time.Hour * 72).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, _ := token.SignedString(jwtSecretKey)
	return &auth.LoginResponse{
		AccessToken: t,
	}, nil
}

func Registration(ctx context.Context, in *auth.RegistrationRequest) (*auth.RegistrationResponse, error) {
	var user *Model.User;

	r := db.DB.Unscoped().Model(&user).First(&user, "email = ?", in.GetEmail())
	if r.RowsAffected > 0 {
		return nil, status.Error(codes.AlreadyExists, "Такой пользователь уже есть")
	}

	pass, _ := auth_action.HashPassword(in.GetPassword())

	newUser := Model.User{
		Email: in.GetEmail(),
		Password: pass,	
		Name: in.GetName(),	
	}

	r = db.DB.Create(&newUser)
	os.Mkdir("files/"+strconv.Itoa(int(user.ID)), os.ModePerm)

	if r.RowsAffected==0{
		return nil, status.Error(codes.Unauthenticated, "Не зарегистрирован")
	}

    payload := jwt.MapClaims{
		"email": user.Email,
        "exp": time.Now().Add(time.Hour * 72).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, _ := token.SignedString(jwtSecretKey)
	
	return &auth.RegistrationResponse{
		AccessToken: t,
	}, nil
}

func DHConnect(ctx context.Context, in *auth.DHConnectRequest) (*auth.DHConnectResponse, error) {
	p,g := helpers.GeneratePubKeys()
	id := uuid.New()
	b, _ := helpers.GeneratePubKey(p, int(g))
	keysUser[id.String()] = KeyUser{
		p: p,
		g: g,
		b: b,
	}
	return &auth.DHConnectResponse{P: p.String(), G: g, UserIdKey: id.String(), B: b.String()}, nil
}

func DHSecondConnect(ctx context.Context, in *auth.DHSecondConnectRequest) (*auth.DHSecondConnectResponse, error) {
	keys := keysUser[in.GetUserIdKey()]
	b, _ := helpers.GeneratePubKey(keys.b, int(keys.g))
	generatedKeys[in.GetUserIdKey()] = b.String()

	return &auth.DHSecondConnectResponse{
		Message: "Ключ успешно создан",
	}, nil
}