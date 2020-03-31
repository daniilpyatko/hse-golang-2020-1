package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func AccessLog(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	res := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Infow("New request",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"time", time.Since(start),
		)
	}
	return http.HandlerFunc(res)
}
