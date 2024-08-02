package helpers

import (
	"mypackages/db"
	Model "mypackages/models"

)

func GetUser(token string) (*Model.User, bool) {
	email:=ParseJWT(token)

	var user Model.User;

	r := db.DB.Model(&user).First(&user, "email=?", email)
	
	if r.RowsAffected==0{

		return nil, false
	}

	return &user, true
}