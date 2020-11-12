package middleware

import (
	"context"
	"net/http"
)

func SetMuxVars(next http.HandlerFunc, key, value string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, key, value)
		next(w, r.WithContext(ctx))
	}
}