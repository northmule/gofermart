package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewGophermartConfig(t *testing.T) {
	testTable := []struct {
		name     string
		env      map[string]string
		args     []string
		expected struct {
			ServerURL   string
			AccrualURL  string
			DatabaseURI string
			LogLevel    string
		}
	}{
		{
			name: "Environment",
			env: map[string]string{
				"RUN_ADDRESS":            "localhost:8081",
				"ACCRUAL_SYSTEM_ADDRESS": "http://localhost:8091",
				"DATABASE_URI":           "postgres://postgres:123@localhost:5456/gofermart?sslmode=disable",
				"LOG_LEVEL":              "info",
			},
			args: []string{},
			expected: struct {
				ServerURL   string
				AccrualURL  string
				DatabaseURI string
				LogLevel    string
			}{
				ServerURL:   "localhost:8081",
				AccrualURL:  "http://localhost:8091",
				DatabaseURI: "postgres://postgres:123@localhost:5456/gofermart?sslmode=disable",
				LogLevel:    "info",
			},
		},
		{
			name: "Flags",
			env: map[string]string{
				"RUN_ADDRESS":            "localhost:8083",
				"ACCRUAL_SYSTEM_ADDRESS": "http://localhost:8093",
				"DATABASE_URI":           "postgres://postgres:123@localhost:5458/gofermart?sslmode=disable",
				"LOG_LEVEL":              "error",
			},
			args: []string{"-a", "localhost:8084", "-d", "postgres://postgres:123@localhost:5459/gofermart?sslmode=disable", "-r", "http://localhost:8094", "-l", "warn"},
			expected: struct {
				ServerURL   string
				AccrualURL  string
				DatabaseURI string
				LogLevel    string
			}{
				ServerURL:   "localhost:8084",
				AccrualURL:  "http://localhost:8094",
				DatabaseURI: "postgres://postgres:123@localhost:5459/gofermart?sslmode=disable",
				LogLevel:    "warn",
			},
		},
		{
			name: "Environment_and_flag",
			env: map[string]string{
				"RUN_ADDRESS":            "localhost:8083",
				"ACCRUAL_SYSTEM_ADDRESS": "http://localhost:8093",
				"DATABASE_URI":           "postgres://postgres:123@localhost:5458/gofermart?sslmode=disable",
				"LOG_LEVEL":              "error",
			},
			args: []string{"-a", "localhost:8084", "-d", "postgres://postgres:123@localhost:5459/gofermart?sslmode=disable", "-r", "http://localhost:8094", "-l", "warn"},
			expected: struct {
				ServerURL   string
				AccrualURL  string
				DatabaseURI string
				LogLevel    string
			}{
				ServerURL:   "localhost:8084",
				AccrualURL:  "http://localhost:8094",
				DatabaseURI: "postgres://postgres:123@localhost:5459/gofermart?sslmode=disable",
				LogLevel:    "warn",
			},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.env {
				t.Setenv(k, v)
			}
			os.Args = append([]string{"test"}, test.args...)

			cfg, err := NewGophermartConfig()
			assert.NoError(t, err)
			assert.Equal(t, test.expected.ServerURL, cfg.ServerURL)
			assert.Equal(t, test.expected.AccrualURL, cfg.AccrualURL)
			assert.Equal(t, test.expected.DatabaseURI, cfg.DatabaseURI)
			assert.Equal(t, test.expected.LogLevel, cfg.LogLevel)
		})
	}
}
