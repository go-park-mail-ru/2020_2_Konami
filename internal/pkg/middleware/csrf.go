package middleware

import (
	"context"
	"konami_backend/logger"
	"konami_backend/proto/csrf"
	"net/http"
)

type CSRFMiddleware struct {
	CsrfClient csrf.CsrfDispatcherClient
	log        *logger.Logger
}

const CSRFValid = "CSRFValid"

func NewCsrfMiddleware(csrfClient csrf.CsrfDispatcherClient, log *logger.Logger) CSRFMiddleware {
	return CSRFMiddleware{
		CsrfClient: csrfClient,
		log:        log,
	}
}

func (m *CSRFMiddleware) CSRFCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authTok, ok := ctx.Value(AuthToken).(string)
		if !ok {
			ctx = context.WithValue(ctx, CSRFValid, false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		CSRFToken := r.Header.Get("Csrf-Token")
		if CSRFToken == "" {
			ctx = context.WithValue(ctx, CSRFValid, false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		isValid, err := m.CsrfClient.Check(context.Background(), &csrf.CsrfToken{Token: CSRFToken, Sid: authTok})
		if err != nil || !isValid.Value {
			ctx = context.WithValue(ctx, CSRFValid, false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx = context.WithValue(ctx, CSRFValid, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
