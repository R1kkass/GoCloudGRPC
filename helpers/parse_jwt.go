package helpers

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey, _ = os.LookupEnv("SECRET_KEY")

func ParseJWT(token string) (string, error) {
	claims := jwt.MapClaims{}
	_, err:= jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return err.Error(), nil
	}

	return claims["email"].(string), nil
}