/*
Package auth deals with generating and checking JWTs
*/
package auth

import (
	"errors"

	"github.com/antonivlev/gql-server/models"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"role":   "user",
	})
	tokenString, errToken := token.SignedString([]byte("verysecret"))
	if errToken != nil {
		return "", errors.New("jwt error: " + errToken.Error())
	}
	return tokenString, nil
}
