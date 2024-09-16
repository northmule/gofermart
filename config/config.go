package config

import (
	"flag"
	"github.com/caarlos0/env"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"os"
	"strings"
)

const DataBaseConnectionTimeOut = 10000
const serverURLDefault = ":8081"
const databaseURIDefault = "postgres://postgres:123@localhost:5456/gofermart?sslmode=disable"
const accrualURLDefault = "http://localhost:8091"

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
	err := instance.env()
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
	serverURL := cf.String("a", serverURLDefault, "адрес и порт запуска сервиса")
	databaseURI := cf.String("d", databaseURIDefault, "адрес подключения к базе данных")
	accrualURL := cf.String("r", accrualURLDefault, "адрес системы расчёта начислений")

	err := cf.Parse(os.Args[1:])
	if err != nil {
		logger.LogSugar.Errorf("флаги конфигурации не разобраны: %s", err)
		return err
	}

	flagsSet := make(map[string]bool)
	cf.Visit(func(f *flag.Flag) {
		flagsSet[f.Name] = true
	})
	var ok bool
	if _, ok = flagsSet["a"]; ok {
		c.ServerURL = *serverURL
	}
	if _, ok = flagsSet["d"]; ok {
		c.DatabaseURI = *databaseURI
	}
	if _, ok = flagsSet["r"]; ok {
		c.AccrualURL = *accrualURL
	}
	// Установка по умолчанию при отсутвии переданных значений
	if c.ServerURL == "" {
		c.ServerURL = *serverURL
	}
	if c.DatabaseURI == "" {
		c.DatabaseURI = *databaseURI
	}
	if c.AccrualURL == "" {
		c.AccrualURL = *accrualURL
	}

	c.ServerURL = strings.ReplaceAll(c.ServerURL, "\"", "")
	c.ServerURL = strings.ReplaceAll(c.ServerURL, " ", "")

	c.DatabaseURI = strings.ReplaceAll(c.DatabaseURI, "\"", "")
	c.DatabaseURI = strings.ReplaceAll(c.DatabaseURI, " ", "")

	c.AccrualURL = strings.ReplaceAll(c.AccrualURL, "\"", "")
	c.AccrualURL = strings.ReplaceAll(c.AccrualURL, " ", "")
	return nil
}
