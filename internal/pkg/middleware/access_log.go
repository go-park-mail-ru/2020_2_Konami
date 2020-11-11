package middleware

import (
	"konami_backend/logger"
	"net/http"
	"time"
)

type AccessLogMiddleware struct {
	Logger *logger.Logger
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func NewAccessLogMiddleware(Logger *logger.Logger) AccessLogMiddleware {
	return AccessLogMiddleware{Logger: Logger}
}

func (m AccessLogMiddleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := statusRecorder{w, 200}
		next.ServeHTTP(&rec, r)
		m.Logger.LogAccess(r, rec.status, time.Since(start))
	})
}
