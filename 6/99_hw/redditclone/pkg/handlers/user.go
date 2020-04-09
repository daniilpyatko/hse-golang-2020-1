package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"redditclone/pkg/session"
	"redditclone/pkg/user"

	"github.com/dgrijalva/jwt-go"
)

// mockgen -source=user.go -destination=user_mock.go -package=handlers UserRepoInterface

type UserRepo interface {
	Authorize(username, password string) (*user.User, error)
	NewUser(username, password string) (*user.User, error)
	GetUserById(Id string) (*user.User, error)
}

type SessionManager interface {
	Create(userId string) (*session.Session, error)
	Check(sessionId string) (bool, error)
	GetUserIdBySessionId(sessionId string) (string, error)
}

type UserHandler struct {
	UserRepo       UserRepo
	SessionManager SessionManager
}

type Claims struct {
	User      *user.User `json:"user"`
	SessionId string     `json:"sessionid"`
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
	}
	curSess, _ := u.SessionManager.Create(curUser.Id)
	u.writeToken(curSess, w)
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var curQuery JSONQuery
	read, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(read, &curQuery)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	username := curQuery.Username
	password := curQuery.Password
	curUser, err := u.UserRepo.NewUser(username, password)
	if err != nil {
		fmt.Println(err)
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
	curSess, _ := u.SessionManager.Create(curUser.Id)
	u.writeToken(curSess, w)
}

func (u *UserHandler) writeToken(curSession *session.Session, w http.ResponseWriter) {
	curUser, _ := u.UserRepo.GetUserById(curSession.UserId)
	tokenString, _ := session.ToToken(curUser, curSession.Id)
	finToken := Token{
		Token: tokenString,
	}
	res, _ := json.Marshal(finToken)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
