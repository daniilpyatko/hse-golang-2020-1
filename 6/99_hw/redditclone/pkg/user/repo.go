package user

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
)

type UserRepo struct {
	data map[string]*User
	mu   *sync.RWMutex
}

func NewRepo() *UserRepo {
	return &UserRepo{
		data: make(map[string]*User),
		mu:   &sync.RWMutex{},
	}
}

var (
	ErrNoUser      = errors.New("User not found")
	ErrBadPassword = errors.New("Invalid Password")
)

func (u *UserRepo) Authorize(username, password string) (*User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, curUser := range u.data {
		if curUser.Username == username {
			if curUser.password == password {
				return curUser, nil
			} else {
				return nil, ErrBadPassword
			}
		}
	}
	return nil, ErrNoUser
}

var (
	ErrUserRegistered = errors.New("User with this username is already registered")
)

func (u *UserRepo) NewUser(username, password string) (*User, error) {
	u.mu.Lock()
	for _, curUser := range u.data {
		if curUser.Username == username {
			return nil, ErrUserRegistered
		}
	}
	rnd := make([]byte, 16)
	rand.Read(rnd)
	randId := base64.URLEncoding.EncodeToString(rnd)
	u.data[randId] = &User{
		Username: username,
		password: password,
		Id:       randId,
	}
	u.mu.Unlock()
	return u.data[randId], nil
}

var (
	ErrUserNotFound = errors.New("User with this Id was not found")
)

func (u *UserRepo) GetUserById(Id string) (*User, error) {
	curUser, ok := u.data[Id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return curUser, nil
}
