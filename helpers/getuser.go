package helpers

import (
	"mypackages/db"
	Model "mypackages/models"
	"strings"
)

func GetUser(jwtToken []string) (*Model.User, bool) {
	jwtToken = strings.Split(jwtToken[0], " ")
	email:=ParseJWT(jwtToken[1])

	var user Model.User;

	r := db.DB.Model(&user).First(&user, "email=?", email)
	
	if r.RowsAffected==0{

		return nil, false
	}

	return &user, true
}