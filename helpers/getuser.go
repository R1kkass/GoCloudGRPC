package helpers

import (
	"errors"
	"mypackages/db"
	Model "mypackages/models"
	"strings"
)

func GetUser(jwtToken []string) (*Model.User, error) {
	if len(jwtToken) == 0 {
		return nil,  errors.New("пользователь не найден")
	}

	jwtToken = strings.Split(jwtToken[0], " ")
	email := ParseJWT(jwtToken[1])

	var user Model.User;
	r := db.DB.Model(&Model.User{}).Where("email=?", email).First(&user)
	
	if r.RowsAffected==0{
		return nil,  errors.New("пользователь не найден")
	}

	return &user, nil
}