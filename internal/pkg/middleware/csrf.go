package middleware

import (
	"context"
	"konami_backend/internal/pkg/csrf"
	"konami_backend/logger"
	"net/http"
)

type CSRFMiddleware struct {
	CsrfUC csrf.UseCase
	log    *logger.Logger
}

const CSRFValid = "CSRFValid"

func NewCsrfMiddleware(csrfUC csrf.UseCase, log *logger.Logger) CSRFMiddleware {
	return CSRFMiddleware{
		CsrfUC: csrfUC,
		log:    log,
	}
}

func (m *CSRFMiddleware) CSRFCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authTok, ok := ctx.Value(AuthToken).(string)
		if !ok {
			m.log.LogWarning("middleware", "CSRFCheck", "context error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		CSRFToken := r.Header.Get("Csrf-Token")
		ok, err := m.CsrfUC.Check(authTok, CSRFToken)
		if err != nil || !ok {
			ctx = context.WithValue(ctx, CSRFValid, false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx = context.WithValue(ctx, CSRFValid, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
