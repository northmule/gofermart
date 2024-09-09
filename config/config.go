package config

import (
	"flag"
	"github.com/caarlos0/env"
	"github.com/northmule/gofermart/internal/app/services/logger"
	"os"
	"strings"
)

const DataBaseConnectionTimeOut = 10000
const ServerURLDefault = ":8081"
const DatabaseURIDefault = "postgres://postgres:123@localhost:5456/gofermart?sslmode=disable"
const AccrualURLDefault = "http://localhost:8091"

type GophermartConfig struct {
	// Адрес сервера и порт
	ServerURL string `env:"RUN_ADDRESS"`
	// Внешняя система расчёта бонусов
	AccrualURL string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	// Строка подключения к БД
	DatabaseURI string `env:"DATABASE_URI"`
}

func NewGophermartConfig() (*GophermartConfig, error) {
	instance := &GophermartConfig{}
	err := instance.flag()
	if err != nil {
		return nil, err
	}

	return instance, instance.flag()
}

func (c *GophermartConfig) env() error {
	err := env.Parse(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *GophermartConfig) flag() error {
	cf := flag.FlagSet{}
	serverURL := cf.String("a", ServerURLDefault, "адрес и порт запуска сервиса")
	databaseURI := cf.String("d", DatabaseURIDefault, "адрес подключения к базе данных")
	accrualURL := cf.String("r", AccrualURLDefault, "адрес системы расчёта начислений")

	err := cf.Parse(os.Args[1:])
	if err != nil {
		logger.LogSugar.Errorf("флаги конфигурации не разобраны: %s", err)
		return err
	}

	if *serverURL != "" {
		c.ServerURL = *serverURL
	}
	if *databaseURI != "" {
		c.DatabaseURI = *databaseURI
	}
	if *accrualURL != "" {
		c.AccrualURL = *accrualURL
	}
	c.DatabaseURI = strings.ReplaceAll(c.DatabaseURI, "\"", "")
	return nil
}
