package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func Panic(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Panic middleware %v", r.URL.Path)
				logger.Errorf("recovered %v", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
