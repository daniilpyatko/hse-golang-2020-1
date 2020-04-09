package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"redditclone/pkg/session"
	"redditclone/pkg/user"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUserHandlerRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// registered user
	userRepoMock := NewMockUserRepo(ctrl)
	sessionManagerMock := NewMockSessionManager(ctrl)
	service := &UserHandler{
		UserRepo:       userRepoMock,
		SessionManager: sessionManagerMock,
	}
	test := map[string]string{
		"username": "user",
		"password": "pass",
		"userid":   "1",
	}
	userRepoMock.EXPECT().NewUser(test["username"], test["password"]).Return(nil, user.ErrUserRegistered)
	query := map[string]string{
		"username": "user",
		"password": "pass",
	}
	jsonQuery, _ := json.Marshal(query)
	r := httptest.NewRequest("POST", "/api/register", bytes.NewReader(jsonQuery))
	w := httptest.NewRecorder()
	service.Register(w, r)

	// OK
	resultUser := &user.User{
		Username: test["username"],
		Id:       test["userid"],
	}
	resultSession := &session.Session{
		UserId: test["userid"],
		Id:     "2",
	}

	userRepoMock.EXPECT().NewUser(test["username"], test["password"]).Return(resultUser, nil)
	sessionManagerMock.EXPECT().Create(resultUser.Id).Return(resultSession, nil)
	userRepoMock.EXPECT().GetUserById(resultUser.Id).Return(resultUser, nil)
	query = map[string]string{
		"username": "user",
		"password": "pass",
	}
	jsonQuery, _ = json.Marshal(query)
	r = httptest.NewRequest("POST", "/api/register", bytes.NewReader(jsonQuery))
	w = httptest.NewRecorder()
	service.Register(w, r)

	// invalid json
	jsonQuery = []byte("not valid json")
	r = httptest.NewRequest("POST", "/api/login", bytes.NewReader(jsonQuery))
	w = httptest.NewRecorder()
	service.Register(w, r)
}

func TestUserHandlerLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// user not registered
	userRepoMock := NewMockUserRepo(ctrl)
	sessionManagerMock := NewMockSessionManager(ctrl)
	service := &UserHandler{
		UserRepo:       userRepoMock,
		SessionManager: sessionManagerMock,
	}
	test := map[string]string{
		"username": "user",
		"password": "pass",
		"userid":   "1",
	}
	userRepoMock.EXPECT().Authorize(test["username"], test["password"]).Return(nil, user.ErrNoUser)
	query := map[string]string{
		"username": "user",
		"password": "pass",
	}
	resultUser := &user.User{
		Username: test["username"],
		Id:       test["userid"],
	}
	resultSession := &session.Session{
		UserId: test["userid"],
		Id:     "2",
	}
	jsonQuery, _ := json.Marshal(query)
	r := httptest.NewRequest("POST", "/api/login", bytes.NewReader(jsonQuery))
	w := httptest.NewRecorder()
	service.Login(w, r)

	// bad password
	userRepoMock.EXPECT().Authorize(test["username"], test["password"]).Return(nil, user.ErrBadPassword)
	// sessionManagerMock.EXPECT().Create(resultUser.Id).Return(resultSession, nil)
	query = map[string]string{
		"username": "user",
		"password": "pass",
	}
	jsonQuery, _ = json.Marshal(query)
	r = httptest.NewRequest("POST", "/api/login", bytes.NewReader(jsonQuery))
	w = httptest.NewRecorder()
	service.Login(w, r)

	// OK
	userRepoMock.EXPECT().Authorize(test["username"], test["password"]).Return(resultUser, nil)
	sessionManagerMock.EXPECT().Create(resultUser.Id).Return(resultSession, nil)
	userRepoMock.EXPECT().GetUserById(resultUser.Id).Return(resultUser, nil)
	query = map[string]string{
		"username": "user",
		"password": "pass",
	}
	jsonQuery, _ = json.Marshal(query)
	r = httptest.NewRequest("POST", "/api/login", bytes.NewReader(jsonQuery))
	w = httptest.NewRecorder()
	service.Login(w, r)

	// invalid json
	jsonQuery = []byte("not valid json")
	r = httptest.NewRequest("POST", "/api/login", bytes.NewReader(jsonQuery))
	w = httptest.NewRecorder()
	service.Login(w, r)
}

// go test -v -coverprofile=handler.out && go tool cover -html=handler.out -o handler.html && rm handler.out
