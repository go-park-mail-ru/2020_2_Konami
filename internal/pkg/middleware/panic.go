package middleware

import (
	"konami_backend/internal/pkg/utils/http_utils"
	hu "konami_backend/internal/pkg/utils/http_utils"
	"konami_backend/logger"
	"net/http"
)

type PanicMiddleware struct {
	Logger *logger.Logger
}

func NewPanicMiddleware(Logger *logger.Logger) PanicMiddleware {
	return PanicMiddleware{Logger: Logger}
}

func (m *PanicMiddleware) PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				http_utils.WriteError(w, &hu.ErrResponse{RespCode: http.StatusInternalServerError})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
