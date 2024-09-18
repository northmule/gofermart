package client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccrualClient_SendOrderNumber(t *testing.T) {
	_, _ = logger.NewLogger("info")
	ctx := context.Background()

	t.Run("Положительный_ответ", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/orders/12345" {
				t.Errorf("Expected path '/api/orders/12345', got %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
			response := ResponseAccrual{Order: "12345", Status: "OK", Accrual: 100.0}
			json.NewEncoder(w).Encode(response)
		}))
		defer ts.Close()
		client := NewAccrualClient(ts.URL, logger.LogSugar)

		response, err := client.SendOrderNumber(ctx, "12345")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if response.Order != "12345" {
			t.Errorf("Expected order '12345', got %s", response.Order)
		}
		if response.Status != "OK" {
			t.Errorf("Expected status 'OK', got %s", response.Status)
		}
		if response.Accrual != 100.0 {
			t.Errorf("Expected accrual '100.0', got %f", response.Accrual)
		}
	})
	t.Run("Ошибка_сервера", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()
		client := NewAccrualClient(ts.URL, logger.LogSugar)

		var expectedError ErrorInternalServerError
		_, err := client.SendOrderNumber(ctx, "12345")
		if err == nil || errors.Is(err, expectedError) {
			t.Errorf("Expected error %s, got %v", expectedError, err)
		}
	})

	t.Run("Перегружен_запросами", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer ts.Close()
		client := NewAccrualClient(ts.URL, logger.LogSugar)

		var expectedError ErrorTooManyRequests
		_, err := client.SendOrderNumber(ctx, "12345")
		if err == nil || errors.Is(err, expectedError) {
			t.Errorf("Expected error %s, got %v", expectedError, err)
		}
	})

	t.Run("Нет_данных", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()
		client := NewAccrualClient(ts.URL, logger.LogSugar)

		var expectedError ErrorNoContent
		_, err := client.SendOrderNumber(ctx, "12345")
		if err == nil || errors.Is(err, expectedError) {
			t.Errorf("Expected error %s, got %v", expectedError, err)
		}
	})
}
