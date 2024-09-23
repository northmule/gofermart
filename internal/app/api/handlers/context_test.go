package handlers

import (
	"context"
	"net/http"
	"testing"
)

func TestAddCommonContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", "value")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if val, ok := r.Context().Value("key").(string); !ok || val != "value" {
			t.Errorf("Expected context value 'value' but got %v", val)
		}
	})

	middleware := AddCommonContext(ctx)
	handler := middleware(nextHandler)

	req, _ := http.NewRequest("GET", "/test", nil)

	handler.ServeHTTP(nil, req)
}
