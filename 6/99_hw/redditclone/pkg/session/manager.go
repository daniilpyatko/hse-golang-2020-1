package session

import (
	"database/sql"
	"fmt"
)

type SessionManager struct {
	DB *sql.DB
}

func NewSessionManager(db *sql.DB) *SessionManager {
	return &SessionManager{
		DB: db,
	}
}

func (s *SessionManager) Create(userId string) (*Session, error) {
	session, _ := NewSession(userId)
	_, err := s.DB.Exec(
		"INSERT INTO sessions (`id`, `userid`) VALUES (?, ?)",
		session.Id,
		session.UserId,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return session, nil
}

func (s *SessionManager) Check(sessionId string) (bool, error) {
	sess := &Session{}
	err := s.DB.QueryRow("SELECT id, userid FROM sessions WHERE id = ?", sessionId).
		Scan(&sess.Id, &sess.UserId)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (s *SessionManager) GetUserIdBySessionId(sessionId string) (string, error) {
	sess := &Session{}
	err := s.DB.QueryRow("SELECT id, userid FROM sessions WHERE id = ?", sessionId).
		Scan(&sess.Id, &sess.UserId)
	if err != nil {
		return "", err
	} else {
		return sess.UserId, nil
	}
}
