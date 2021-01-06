package handlers

import (
	"errors"
	"time"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/internal/environment"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/models"
	jwt "github.com/form3tech-oss/jwt-go"
)

func genereateJWT(claims models.JWTLogin) (string, error) {
	//Get Key
	secrectKey := environment.AccessENV("SECRECT_KEY")

	if secrectKey == "" {
		return "", errors.New("Error in get Secrect Key From ENV")
	}

	sign := []byte(secrectKey)

	//Make Claims
	payload := jwt.MapClaims{
		"user_id":   claims.ID,
		"user_name": claims.Username,
		"email":     claims.Email,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	//Sign the JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signedToken, err := token.SignedString(sign)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
