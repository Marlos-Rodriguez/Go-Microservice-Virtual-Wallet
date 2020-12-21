package handlers

import (
	"errors"
	"time"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/models"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/environment"
	jwt "github.com/dgrijalva/jwt-go"
)

func genereateJWT(claims models.JWTLogin) (string, error) {
	secrectKey, sucess := environment.AccessENV("SECRECT_KEY")

	if !sucess {
		return "", errors.New("Error in get Secrect Key From ENV")
	}

	sign := []byte(secrectKey)

	payload := jwt.MapClaims{
		"user_id":   claims.ID,
		"user_name": claims.Username,
		"email":     claims.Email,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signedToken, err := token.SignedString(sign)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
