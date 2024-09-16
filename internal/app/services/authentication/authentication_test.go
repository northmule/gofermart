package authentication

import (
	"github.com/northmule/gophermart/internal/app/services/logger"
	"net/http"
	"testing"
	"time"
)

func TestGetUserToken(t *testing.T) {
	_, _ = logger.NewLogger("info")

	req1 := &http.Request{Header: http.Header{}}
	cookie1 := &http.Cookie{Name: CookieAuthName, Value: "valid_token"}
	req1.AddCookie(cookie1)
	if GetUserToken(req1) != "valid_token" {
		t.Errorf("Expected valid_token, got %s", GetUserToken(req1))
	}

	req2 := &http.Request{Header: http.Header{}}
	if GetUserToken(req2) != "" {
		t.Errorf("Expected empty token, got %s", GetUserToken(req2))
	}

	req3 := &http.Request{Header: http.Header{}}
	cookie3 := &http.Cookie{Name: CookieAuthName, Value: "invalid_token"}
	req3.AddCookie(cookie3)
	if GetUserToken(req3) != "invalid_token" {
		t.Errorf("Expected invalid_token, got %s", GetUserToken(req3))
	}
}

func TestGenerateToken(t *testing.T) {
	_, _ = logger.NewLogger("info")

	userUUID := "user123"
	exp := time.Hour * 600
	secretKey := "super_secret_key_gophermart"

	token, _ := GenerateToken(userUUID, exp, secretKey)
	if len(token) == 0 {
		t.Error("Expected non-empty token, got empty")
	}
}

func TestValidateToken(t *testing.T) {
	_, _ = logger.NewLogger("info")

	userUUID := "user123"
	secretKey := "super_secret_key_gophermart"

	token, _ := GenerateToken(userUUID, HMACTokenExp, secretKey)

	if !ValidateToken(userUUID, token, secretKey) {
		t.Error("Expected token to be valid, but it was not")
	}

	invalidToken := "invalid_token"
	if ValidateToken(userUUID, invalidToken, secretKey) {
		t.Error("Expected token to be invalid, but it was valid")
	}
}
