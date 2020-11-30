package middleware

import (
	"context"
	"net/http"
)

type RouteArgs struct {
	Key   string
	Value interface{}
}

func SetMuxVars(next http.HandlerFunc, args []RouteArgs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		for _, val := range args {
			ctx = context.WithValue(ctx, val.Key, val.Value)
		}
		next(w, r.WithContext(ctx))
	}
}