package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"redditclone/pkg/user"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("Dgbaiubfaskjdljf")

type Session struct {
	Id     string
	UserId string
}

type Claims struct {
	User      *user.User `json:"user"`
	SessionId string     `json:"sessionid"`
	jwt.StandardClaims
}

func NewSession(userId string) (*Session, error) {
	rnd := make([]byte, 16)
	rand.Read(rnd)
	id := base64.URLEncoding.EncodeToString(rnd)
	return &Session{
		Id:     id,
		UserId: userId,
	}, nil
}

func ToToken(user *user.User, sessionId string) (string, error) {
	result := &Claims{
		User:      user,
		SessionId: sessionId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, result)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString, nil
}

var (
	ErrInvalidToken = errors.New("Invalid jwt token")
)

func ToSessionId(token string) (string, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		return "", ErrInvalidToken
	}
	return claims.SessionId, nil
}
