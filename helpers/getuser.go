package helpers

import (
	"errors"
	"fmt"
	"mypackages/db"
	Model "mypackages/models"
	"strings"
)

func GetUser(jwtToken []string) (*Model.User, error) {
	jwtToken = strings.Split(jwtToken[0], " ")
	email := ParseJWT(jwtToken[1])
	fmt.Println(email)

	var user Model.User;
	r := db.DB.Model(&Model.User{}).Where("email=?", email).First(&user)
	
	if r.RowsAffected==0{
		return nil,  errors.New("пользователь не найден")
	}

	return &user, nil
}