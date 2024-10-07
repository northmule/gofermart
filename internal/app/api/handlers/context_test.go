package handlers

import (
	"context"
	"github.com/northmule/gophermart/internal/app/api/rctx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAddCommonContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), rctx.UserCtxKey, "value")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(rctx.UserCtxKey).(string)
		assert.Equal(t, "value", val)
	})

	middleware := AddCommonContext(ctx)
	handler := middleware(nextHandler)

	req, _ := http.NewRequest("GET", "/test", nil)

	handler.ServeHTTP(nil, req)
}
