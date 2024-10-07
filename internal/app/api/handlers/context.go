package handlers

import (
	"context"
	"net/http"
)

func AddCommonContext(ctx context.Context) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			rr := r.Clone(ctx)
			next.ServeHTTP(w, rr)
		}
		return http.HandlerFunc(fn)
	}
}
