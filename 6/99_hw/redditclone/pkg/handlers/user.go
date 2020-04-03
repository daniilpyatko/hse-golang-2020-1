package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"redditclone/pkg/user"

	"github.com/dgrijalva/jwt-go"
)

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
	err := json.Unmarshal(read, &curQuery)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
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

func writeToken(curUser *user.User, w http.ResponseWriter) {
	tokenString, _ := user.UserToToken(curUser)
	finToken := Token{
		Token: tokenString,
	}
	res, _ := json.Marshal(finToken)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
