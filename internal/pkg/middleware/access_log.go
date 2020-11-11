package middleware

import (
	"konami_backend/logger"
	"net/http"
	"time"
)

type AccessLogMiddleware struct {
	Logger *logger.Logger
}

func NewAccessLogMiddleware(Logger *logger.Logger) AccessLogMiddleware {
	return AccessLogMiddleware{Logger: Logger}
}

func (m AccessLogMiddleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		m.Logger.LogAccess(r, time.Since(start))
	})
}
