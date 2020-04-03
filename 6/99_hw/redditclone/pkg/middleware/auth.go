package middleware

import (
	"context"
	"net/http"
	"redditclone/pkg/user"
	"strings"

	"go.uber.org/zap"
)

func Auth(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("Auth middleware")
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")
		curUser, err := user.TokenToUser(token)
		if err != nil {
			logger.Errorf("Token parse error %v", err)
			logger.Errorf("Token %v", token)
			return
		}
		ctx := context.WithValue(r.Context(), "Id", curUser.Id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
