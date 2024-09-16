package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/northmule/gofermart/config"
	_ "go.uber.org/mock/mockgen/model"
	"time"
)

const CodeErrorDuplicateKey = "23505"

type DBQuery interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	PingContext(ctx context.Context) error
	Begin() (*sql.Tx, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
<<<<<<< HEAD
	Prepare(query string) (*sql.Stmt, error)
=======
>>>>>>> 94746e2 (базовая структура)
}

type PostgresStorage struct {
	DB    DBQuery
<<<<<<< HEAD
	RawDB *sql.DB
=======
	SqlDB *sql.DB
>>>>>>> 94746e2 (базовая структура)
}

// NewPostgresStorage PostgresStorage настройка подключения к БД
func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	instance := &PostgresStorage{
		DB:    db,
<<<<<<< HEAD
		RawDB: db,
=======
		SqlDB: db,
>>>>>>> 94746e2 (базовая структура)
	}

	return instance, nil
}

func (p *PostgresStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	return p.DB.PingContext(ctx)
}
