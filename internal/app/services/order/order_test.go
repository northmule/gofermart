package order

import (
	"github.com/northmule/gophermart/internal/app/services/logger"
	"testing"
)

func TestNewOrderService(t *testing.T) {
	s := NewOrderService()
	if s == nil {
		t.Error("Expected OrderService to be initialized")
	}
	if s == nil || s.alg != NumberValidateAlg {
		t.Errorf("Expected alg to be %s, but got %s", NumberValidateAlg, s.alg)
	}
}

func TestValidateOrderNumber(t *testing.T) {
	_, _ = logger.NewLogger("info")
	os := NewOrderService()

	if !os.ValidateOrderNumber(49927398716) {
		t.Error("Expected the number to be valid")
	}

	if os.ValidateOrderNumber(1234567812345678) {
		t.Error("Expected the number to be invalid")
	}
}

func TestLuhnValid(t *testing.T) {
	s := NewOrderService()

	if !s.luhnValid(49927398716) {
		t.Error("Expected the number to be valid")
	}

	if s.luhnValid(1234567812345678) {
		t.Error("Expected the number to be invalid")
	}
}

func TestLuhnChecksum(t *testing.T) {
	s := NewOrderService()

	expected := 4
	result := s.luhnChecksum(4992739871)
	if result != expected {
		t.Errorf("Expected checksum to be %d, but got %d", expected, result)
	}
}
