package storage

import (
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	ss := NewSessionStorage()
	token := "testToken"
	expire := time.Now().Add(time.Hour)

	ss.Add(token, expire)

	if _, ok := ss.(*SessionStorage).values[token]; !ok {
		t.Errorf("Token not added to storage")
	}
}

func TestIsValid(t *testing.T) {
	ss := NewSessionStorage()
	token := "testToken"
	expire := time.Now().Add(time.Hour)

	if ss.IsValid(token) {
		t.Errorf("Token should be invalid initially")
	}

	ss.Add(token, expire)

	if !ss.IsValid(token) {
		t.Errorf("Token should be valid after adding")
	}

	time.Sleep(time.Hour)
	if ss.IsValid(token) {
		t.Errorf("Token should be invalid after expiration")
	}
}
