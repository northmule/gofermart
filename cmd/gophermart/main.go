package main

import (
	"context"
	"errors"
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
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Запуск приложения Gophermart")
	appCtx, appStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appStop()
	if err := run(appCtx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	loggerService, err := logger.NewLogger("info")
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
	store, err := storage.NewPostgresStorage(cfg.DatabaseURI, ctx)
	if err != nil {
		return err
	}
	logger.LogSugar.Info("Проверка подключения к БД")
	err = store.Ping()
	if err != nil {
		return err
	}
	logger.LogSugar.Info("Инициализация миграций")
	migrations := db.NewMigrations(store.RawDB, ctx)
	err = migrations.Up()
	if err != nil {
		return err
	}
	logger.LogSugar.Info("Инициализация менеджера репозитариев")
	repositoryManager := repository.NewManager(store.DB, ctx)
	routes := api.NewAppRoutes(repositoryManager, ctx, loggerService)

	logger.LogSugar.Info("Инициализация клиента Accrual")
	accrualClient := client.NewAccrualClient(cfg.AccrualURL, logger.LogSugar, ctx)
	logger.LogSugar.Info("Инициализация worker-ов")
	_ = job.NewWorker(repositoryManager, accrualClient, ctx)

	httpServer := http.Server{
		Addr:    cfg.ServerURL,
		Handler: routes,
	}
	go func() {
		logger.LogSugar.Infof("Running server on - %s", cfg.ServerURL)
		err = httpServer.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.LogSugar.Fatal(err)
		}
		if errors.Is(err, http.ErrServerClosed) {
			logger.LogSugar.Info("Сервер остановлен")
		}
	}()

	<-ctx.Done()
	logger.LogSugar.Info("Получин сигнал. Останавливаю сервер...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = httpServer.Shutdown(shutdownCtx)
	if err != nil {
		return err
	}

	defer cancel()

	return nil
}
