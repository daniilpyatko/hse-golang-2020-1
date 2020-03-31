package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"redditclone/pkg/user"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("Dgbaiubfaskjdljf")

type UserHandler struct {
	UserRepo *user.UserRepo
}

type Claims struct {
	User *user.User `json:"user"`
	jwt.StandardClaims
}

type Token struct {
	Token string `json:"token"`
}

type JSONQuery struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var curQuery JSONQuery
	read, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(read, &curQuery)
	username := curQuery.Username
	password := curQuery.Password
	curUser, err := u.UserRepo.Authorize(username, password)
	if err != nil {
		if err == user.ErrNoUser {
			message := map[string]string{
				"message": "user not found",
			}
			jsonmessage, _ := json.Marshal(message)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(jsonmessage)
			return
		} else if err == user.ErrBadPassword {
			message := map[string]string{
				"message": "invalid password",
			}
			jsonmessage, _ := json.Marshal(message)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(jsonmessage)
			return
		}
		return
	}
	writeToken(curUser, w)
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var curQuery JSONQuery
	read, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(read, &curQuery)
	username := curQuery.Username
	password := curQuery.Password
	curUser, err := u.UserRepo.NewUser(username, password)
	if err != nil {
		message := map[string][]map[string]string{
			"errors": []map[string]string{
				map[string]string{
					"location": "body",
					"param":    "username",
					"value":    username,
					"msg":      "already exists",
				},
			},
		}
		jsonmessage, _ := json.Marshal(message)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(jsonmessage)
		return
	}
	writeToken(curUser, w)
}

func writeToken(user *user.User, w http.ResponseWriter) {
	result := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, result)
	tokenString, _ := token.SignedString(jwtKey)
	finToken := Token{
		Token: tokenString,
	}
	res, _ := json.Marshal(finToken)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

var (
	ErrInvalidToken = errors.New("Invalid jwt token")
)

func TokenToUser(token string) (*user.User, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		return nil, ErrInvalidToken
	}
	return claims.User, nil
}
