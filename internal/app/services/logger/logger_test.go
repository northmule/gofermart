package logger

import (
	"net/http"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	_, err := NewLogger("info")
	if err != nil {
		t.Errorf("Expected logger creation: %v", err)
	}
}

func TestLogger_Print(t *testing.T) {
	logger, _ := NewLogger("info")
	logger.Print("Test message")
}

func TestLogger_NewLogEntry(t *testing.T) {
	logger, _ := NewLogger("info")
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	entry := logger.NewLogEntry(req)
	if entry == nil {
		t.Error("Expected LogEntry creation")
	}
}

func TestLogEntry_Write(t *testing.T) {
	logger, _ := NewLogger("info")
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	entry := logger.NewLogEntry(req)
	entry.Write(200, 100, http.Header{}, 100*time.Millisecond, "additional information")
}

func TestLogEntry_Panic(t *testing.T) {
	logger, _ := NewLogger("info")
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	entry := logger.NewLogEntry(req)
	entry.Panic("Test panic", []byte("test trace"))
}
