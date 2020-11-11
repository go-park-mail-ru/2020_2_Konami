package middleware

import (
	"context"
	"konami_backend/internal/pkg/profile"
	"konami_backend/internal/pkg/session"
	"net/http"
)

const (
	AuthStatus = "authStatus"
	UserID     = "userId"
	AuthToken  = "authToken"
)

type AuthMiddleware struct {
	ProfileUC profile.UseCase
	SessionUC session.UseCase
}

func NewAuthMiddleware(ProfileUC profile.UseCase, SessionUC session.UseCase) AuthMiddleware {
	return AuthMiddleware{
		ProfileUC: ProfileUC,
		SessionUC: SessionUC,
	}
}

func (am *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token, err := r.Cookie("authToken")
		if err != nil {
			ctx = context.WithValue(ctx, AuthStatus, true)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		uId, err := am.SessionUC.GetUserId(token.Value)
		if err != nil {
			ctx = context.WithValue(ctx, AuthStatus, true)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx = context.WithValue(ctx, AuthStatus, true)
		ctx = context.WithValue(ctx, UserID, uId)
		ctx = context.WithValue(ctx, AuthToken, token.Value)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
