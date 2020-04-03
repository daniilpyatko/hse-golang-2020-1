package user

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Username string `json:"username"`
	Id       string `json:"id"`
	password string
}

type Claims struct {
	User *User `json:"user"`
	jwt.StandardClaims
}

var jwtKey = []byte("Dgbaiubfaskjdljf")

var (
	ErrInvalidToken = errors.New("Invalid jwt token")
)

func TokenToUser(token string) (*User, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		return nil, ErrInvalidToken
	}
	return claims.User, nil
}

func UserToToken(user *User) (string, error) {
	result := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, result)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString, nil
}
