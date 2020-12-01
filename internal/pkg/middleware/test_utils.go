package middleware

import (
	"context"
	"net/http"
)

type RouteArgs struct {
	Key   string
	Value interface{}
}

type QueryArgs struct {
	Key   string
	Value string
}

func SetVars(next http.HandlerFunc, args []QueryArgs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		for _, val := range args {
			q.Add(val.Key, val.Value)

		}
		r.URL.RawQuery = q.Encode()
		next(w, r)
	}
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