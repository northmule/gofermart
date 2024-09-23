package handlers

import (
	"context"
	"github.com/northmule/gophermart/internal/app/api/rctx"
	"net/http"
	"testing"
)

func TestAddCommonContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), rctx.UserCtxKey, "value")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if val, ok := r.Context().Value(rctx.UserCtxKey).(string); !ok || val != "value" {
			t.Errorf("Expected context value 'value' but got %v", val)
		}
	})

	middleware := AddCommonContext(ctx)
	handler := middleware(nextHandler)

	req, _ := http.NewRequest("GET", "/test", nil)

	handler.ServeHTTP(nil, req)
}
