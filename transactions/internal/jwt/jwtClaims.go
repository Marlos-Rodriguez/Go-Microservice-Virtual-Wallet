package jwt

import (
	jwt "github.com/form3tech-oss/jwt-go"
)

var (
	//UserID the ID of User
	UserID string
	//Username The Username
	Username string
	//UserEmail email of User
	UserEmail string
)

//GetClaims Get the values from JWT
func GetClaims(tk jwt.Token) error {
	//Assing the claims in a variable
	claims := tk.Claims.(jwt.MapClaims)

	UserID = claims["user_id"].(string)
	Username = claims["user_name"].(string)
	UserEmail = claims["email"].(string)

	return nil
}
