/*
Package auth deals with generating and checking JWTs
*/
package auth

import (
	"errors"
	"fmt"

	"github.com/antonivlev/gql-server/models"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID": user.ID,
	})
	tokenString, errToken := token.SignedString([]byte("verysecret"))
	if errToken != nil {
		return "", errors.New("jwt error: " + errToken.Error())
	}
	return tokenString, nil
}

func GetUserIDFromToken(tokenString string) string {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("verysecret"), nil
	})
	if err != nil {
		fmt.Println("GetUserFromToken error: ", err.Error())
		return ""
	}
	userID, ok := token.Claims.(jwt.MapClaims)["ID"].(string)
	if !ok {
		fmt.Println("GetUserFromToken error: type conversion in claims")
		return ""
	}
	return userID
}
