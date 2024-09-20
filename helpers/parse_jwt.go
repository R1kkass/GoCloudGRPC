package helpers

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey, _ = os.LookupEnv("SECRET_KEY")

func ParseJWT(token string) string {
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	return claims["email"].(string)
}