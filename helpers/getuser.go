package helpers

import (
	"errors"
	"strings"

	"github.com/R1kkass/GoCloudGRPC/db"
	Model "github.com/R1kkass/GoCloudGRPC/models"
)

func GetUser(jwtToken []string) (*Model.User, error) {
	if len(jwtToken) == 0 {
		return nil, errors.New("пользователь не найден")
	}

	jwtToken = strings.Split(jwtToken[0], " ")
	email, err := ParseJWT(jwtToken[1])

	if err != nil {
		return nil, errors.New("токен не валиден")
	}

	var user Model.User
	r := db.DB.Model(&Model.User{}).Where("email=?", email).First(&user)

	if r.RowsAffected == 0 {
		return nil, errors.New("пользователь не найден")
	}

	return &user, nil
}
