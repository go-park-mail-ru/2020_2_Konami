package middleware

import (
	"context"
	"konami_backend/internal/pkg/profile"
	"konami_backend/proto/auth"
	"net/http"
)

const (
	AuthStatus = "authStatus"
	UserID     = "userId"
	AuthToken  = "authToken"
)

type AuthMiddleware struct {
	ProfileUC   profile.UseCase
	AuthChecker auth.AuthCheckerClient
}

func NewAuthMiddleware(ProfileUC profile.UseCase, AuthChecker auth.AuthCheckerClient) AuthMiddleware {
	return AuthMiddleware{
		ProfileUC:   ProfileUC,
		AuthChecker: AuthChecker,
	}
}

func (am *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token, err := r.Cookie("authToken")
		if err != nil {
			ctx = context.WithValue(ctx, AuthStatus, false) // nolint:staticcheck
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		sess, err := am.AuthChecker.Check(context.Background(), &auth.SessionToken{Token: token.Value})
		if err != nil {
			ctx = context.WithValue(ctx, AuthStatus, false) // nolint:staticcheck
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx = context.WithValue(ctx, AuthStatus, true)         // nolint:staticcheck
		ctx = context.WithValue(ctx, UserID, int(sess.UserId)) // nolint:staticcheck
		ctx = context.WithValue(ctx, AuthToken, token.Value)   // nolint:staticcheck

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
