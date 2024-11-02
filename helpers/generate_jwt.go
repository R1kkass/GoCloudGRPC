package helpers

import (
	"time"

	Model "github.com/R1kkass/GoCloudGRPC/models"

	"github.com/golang-jwt/jwt"
)

var jwtSecretKey = []byte(secretKey)

func GenerateJWTToken(user *Model.User, hashKey string) (*string, error) {
	payload := jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(jwtSecretKey)
	secretToken := Encrypt(t, hashKey)
	if err != nil {
		return nil, err
	}

	return &secretToken, nil
}
