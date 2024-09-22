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
	fmt.Println("Инициализация конфигурации")
	cfg, err := config.NewGophermartConfig()
	if err != nil {
		return err
	}
	_, err = logger.NewLogger(cfg.LogLevel)
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
	err = store.Ping(ctx)
	if err != nil {
		return err
	}
	logger.LogSugar.Info("Инициализация миграций")
	migrations := db.NewMigrations(store.RawDB)
	err = migrations.Up(ctx)
	if err != nil {
		return err
	}
	logger.LogSugar.Info("Инициализация менеджера репозитариев")
	repositoryManager := repository.NewManager(store.DB)

	logger.LogSugar.Info("Инициализация клиента Accrual")
	accrualClient := client.NewAccrualClient(cfg.AccrualURL, logger.LogSugar)
	logger.LogSugar.Info("Инициализация worker-ов")
	worker := job.NewWorker(repositoryManager, accrualClient)
	worker.Run(ctx)

	logger.LogSugar.Info("Подготовка сервера к запуску")
	routes := api.NewAppRoutes(repositoryManager, store.DB)
	httpServer := http.Server{
		Addr:    cfg.ServerURL,
		Handler: routes.DefiningAppRoutes(ctx),
	}
	go func() {
		<-ctx.Done()
		logger.LogSugar.Info("Получин сигнал. Останавливаю сервер...")

		shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()
		err = httpServer.Shutdown(shutdownCtx)
		if err != nil {
			logger.LogSugar.Error(err)
		}
	}()

	logger.LogSugar.Infof("Running server on - %s", cfg.ServerURL)
	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	if errors.Is(err, http.ErrServerClosed) {
		logger.LogSugar.Info("Сервер остановлен")
	}

	return nil
}
