package config

import (
	"os"
	"testing"
)

func TestNewGophermartConfig(t *testing.T) {

	t.Setenv("RUN_ADDRESS", "localhost:8081")
	t.Setenv("ACCRUAL_SYSTEM_ADDRESS", "http://localhost:8091")
	t.Setenv("DATABASE_URI", "postgres://postgres:123@localhost:5456/gofermart?sslmode=disable")

	os.Args = []string{"test"}
	config, err := NewGophermartConfig()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Проверяем значения конфигурации
	if config.ServerURL != "localhost:8081" {
		t.Errorf("Expected ServerURL to be 'localhost:8081', but got '%s'", config.ServerURL)
	}
	if config.AccrualURL != "http://localhost:8091" {
		t.Errorf("Expected AccrualURL to be 'http://localhost:8091', but got '%s'", config.AccrualURL)
	}
	if config.DatabaseURI != "postgres://postgres:123@localhost:5456/gofermart?sslmode=disable" {
		t.Errorf("Expected DatabaseURI to be 'postgres://postgres:123@localhost:5456/gofermart?sslmode=disable', but got '%s'", config.DatabaseURI)
	}
}

func TestNewGophermartConfigWithFlags(t *testing.T) {
	// Устанавливаем аргументы командной строки для теста
	os.Args = []string{"cmd", "-a", "localhost:8081", "-d", "postgres"}

	// Вызываем функцию NewGophermartConfig
	config, err := NewGophermartConfig()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Проверяем значения конфигурации
	if config.ServerURL != "localhost:8081" {
		t.Errorf("Expected ServerURL to be 'localhost:8081', but got '%s'", config.ServerURL)
	}
	if config.AccrualURL != "http://localhost:8091" {
		t.Errorf("Expected AccrualURL to be 'http://localhost:8091', but got '%s'", config.AccrualURL)
	}
	if config.DatabaseURI != "postgres" {
		t.Errorf("Expected DatabaseURI to be 'postgres', but got '%s'", config.DatabaseURI)
	}
}
func TestNewGophermartConfigWithMixedSources(t *testing.T) {
	t.Setenv("RUN_ADDRESS", "env_server")
	t.Setenv("ACCRUAL_SYSTEM_ADDRESS", "env_accrual")
	t.Setenv("DATABASE_URI", "env_database")

	os.Args = []string{"cmd", "-a", "localhost:8081", "-d", "postgres://postgres:123@localhost:5456/gofermart?sslmode=disable", "-r", "http://localhost:8091"}

	config, err := NewGophermartConfig()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Проверяем значения конфигурации
	if config.ServerURL != "localhost:8081" {
		t.Errorf("Expected ServerURL to be 'localhost:8081', but got '%s'", config.ServerURL)
	}
	if config.AccrualURL != "http://localhost:8091" {
		t.Errorf("Expected AccrualURL to be 'http://localhost:8091', but got '%s'", config.AccrualURL)
	}
	if config.DatabaseURI != "postgres://postgres:123@localhost:5456/gofermart?sslmode=disable" {
		t.Errorf("Expected DatabaseURI to be 'postgres://postgres:123@localhost:5456/gofermart?sslmode=disable', but got '%s'", config.DatabaseURI)
	}
}
