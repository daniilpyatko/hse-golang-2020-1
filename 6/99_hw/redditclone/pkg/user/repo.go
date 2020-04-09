package user

import (
	"database/sql"
	"errors"
	"fmt"
	"redditclone/pkg/random"
	"sync"
)

type UserRepo struct {
	DB *sql.DB
	mu *sync.RWMutex
	rn *random.Generator
}

func NewRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		DB: db,
		mu: &sync.RWMutex{},
		rn: random.NewGenerator(true),
	}
}

var (
	ErrNoUser      = errors.New("User not found")
	ErrBadPassword = errors.New("Invalid Password")
)

func (u *UserRepo) Authorize(username, password string) (*User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	curUser := User{}
	err := u.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).
		Scan(&curUser.Id, &curUser.Username, &curUser.password)
	if err != nil {
		return nil, ErrNoUser
	}
	if curUser.password != password {
		return nil, ErrBadPassword
	}
	return &curUser, nil
}

var (
	ErrUserRegistered = errors.New("User with this username is already registered")
)

func (u *UserRepo) NewUser(username, password string) (*User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	curUser := User{}
	err := u.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).
		Scan(&curUser.Id, &curUser.Username, &curUser.password)
	if err == nil {
		return nil, ErrUserRegistered
	} else {
		randId := u.rn.GetString()
		newUser := &User{
			Username: username,
			password: password,
			Id:       randId,
		}
		fmt.Println(newUser)
		u.DB.Exec(
			"INSERT INTO users (`id`, `username`, `password`) VALUES (?, ?, ?)",
			newUser.Id,
			newUser.Username,
			newUser.password,
		)
		return newUser, nil
	}
}

var (
	ErrUserNotFound = errors.New("User with this Id was not found")
)

func (u *UserRepo) GetUserById(Id string) (*User, error) {
	curUser := User{}
	err := u.DB.QueryRow("SELECT id, username, password FROM users WHERE id = ?", Id).
		Scan(&curUser.Id, &curUser.Username, &curUser.password)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &curUser, nil
}
