package middleware

import (
	"context"
	"net/http"
	"redditclone/pkg/session"
	"strings"

	"go.uber.org/zap"
)

func Auth(s *session.SessionManager, logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("Auth middleware")
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")
		// curUser, err := user.TokenToUser(token)
		curSess, err := session.ToSessionId(token)
		if ok, _ := s.Check(curSess); !ok {
			logger.Errorf("Session doesn't exist")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err != nil {
			logger.Errorf("Token parse error %v", err)
			logger.Errorf("Token %v", token)
			return
		}
		curUserId, _ := s.GetUserIdBySessionId(curSess)
		ctx := context.WithValue(r.Context(), "Id", curUserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
