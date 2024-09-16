package db

import (
	"context"
	"database/sql"
	"embed"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"github.com/pressly/goose/v3"
)

type Migrations struct {
	mFS   embed.FS
	sqlDB *sql.DB
	ctx   context.Context
}

//go:embed migrations/*.sql
var migrationsFS embed.FS

func NewMigrations(db *sql.DB, ctx context.Context) *Migrations {
	instance := Migrations{}
	instance.mFS = migrationsFS
	instance.sqlDB = db
	instance.ctx = ctx
	return &instance
}

func (m *Migrations) Up() error {
	logger.LogSugar.Info("Запуск миграции")
	goose.SetBaseFS(m.mFS)
	if err := goose.SetDialect("postgres"); err != nil {
		logger.LogSugar.Error(err)
		return err
	}

	if err := goose.UpContext(m.ctx, m.sqlDB, "migrations"); err != nil {
		logger.LogSugar.Error(err)
		return err
	}

	return nil
}
