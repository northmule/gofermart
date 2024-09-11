package storage

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/northmule/gophermart/config"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"time"
)

const CodeErrorDuplicateKey = "23505"

type DBQuery interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	PingContext(ctx context.Context) error
	Begin() (*sql.Tx, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
}

type TxDBQuery interface {
	TxQueryRowContext(ctx context.Context, query string, args ...any) (*sql.Row, error)
	TxOpen() error
	TxRollback() error
	TxCommit() error
}

type PostgresStorage struct {
	DB    DBQuery
	RawDB *sql.DB
	TxDB  TxDBQuery
	tx    *sql.Tx
	ctx   context.Context
}

// NewPostgresStorage PostgresStorage настройка подключения к БД
func NewPostgresStorage(dsn string, ctx context.Context) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	instance := &PostgresStorage{
		DB:    db,
		RawDB: db,
		ctx:   ctx,
	}

	return instance, nil
}

func (p *PostgresStorage) Ping() error {
	ctx, cancel := context.WithTimeout(p.ctx, config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	return p.DB.PingContext(ctx)
}

func (p *PostgresStorage) TxQueryRowContext(query string, args ...any) (*sql.Row, error) {
	rows := p.tx.QueryRowContext(p.ctx, query, args)
	err := rows.Err()
	if err != nil {
		err = errors.Join(err, p.tx.Rollback())
		logger.LogSugar.Error(err)
		return nil, err
	}

	return rows, nil

}

func (p *PostgresStorage) TxOpen() error {
	var err error
	p.tx, err = p.DB.Begin()
	if err != nil {
		logger.LogSugar.Error(err)
		return err
	}
	return nil
}

func (p *PostgresStorage) TxRollback() error {
	return p.tx.Rollback()
}
func (p *PostgresStorage) TxCommit() error {
	return p.tx.Commit()
}
