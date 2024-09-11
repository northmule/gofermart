package main

import (
	"fmt"
	"github.com/northmule/gophermart/config"
	"github.com/northmule/gophermart/db"
	"github.com/northmule/gophermart/internal/accrual/api/client"
	"github.com/northmule/gophermart/internal/app/api"
	"github.com/northmule/gophermart/internal/app/repository"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"github.com/northmule/gophermart/internal/app/storage"
	job "github.com/northmule/gophermart/internal/app/worker"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Запуск приложения Gophermart")
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	err := logger.NewLogger("info")
	if err != nil {
		return err
	}
	logger.LogSugar.Info("Инициализация конфигурации")
	cfg, err := config.NewGophermartConfig()
	if err != nil {
		return err
	}
	logger.LogSugar.Infof("Конфигурация приложения %#v", cfg)
	logger.LogSugar.Info("Инициализация базы данных")
	store, err := storage.NewPostgresStorage(cfg.DatabaseURI)
	if err != nil {
		return err
	}
	logger.LogSugar.Info("Проверка подключения к БД")
	err = store.Ping()
	if err != nil {
		return err
	}
	logger.LogSugar.Info("Инициализация миграций")
	migrations := db.NewMigrations(store.RawDB)
	err = migrations.Up()
	if err != nil {
		return err
	}
	logger.LogSugar.Info("Инициализация менеджера репозитариев")
	repositoryManager := repository.NewManager(store.DB)
	routes := api.NewAppRoutes(repositoryManager)

	logger.LogSugar.Info("Инициализация клиента Accrual")
	accrualClient := client.NewAccrualClient(cfg.AccrualURL, logger.LogSugar)
	logger.LogSugar.Info("Инициализация worker-ов")
	_ = job.NewWorker(repositoryManager, accrualClient)

	logger.LogSugar.Infof("Running server on - %s", cfg.ServerURL)
	return http.ListenAndServe(cfg.ServerURL, routes)
}
