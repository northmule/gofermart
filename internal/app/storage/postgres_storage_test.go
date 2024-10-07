package storage

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"github.com/stretchr/testify/suite"
)

type PostgresStorageTestSuite struct {
	suite.Suite
	db      *sql.DB
	mock    sqlmock.Sqlmock
	storage *PostgresStorage
	ctx     context.Context
	dsn     string
}

func (s *PostgresStorageTestSuite) SetupTest() {
	_, _ = logger.NewLogger("info")
	var err error
	s.ctx = context.Background()
	s.dsn = "user=test dbname=test sslmode=disable"

	s.db, s.mock, err = sqlmock.New()
	s.Require().NoError(err)

	instance, _ := NewPostgresStorage(s.dsn)
	instance.DB = s.db
	s.storage = instance
}

func (s *PostgresStorageTestSuite) TearDownTest() {
	s.db.Close()
}

func (s *PostgresStorageTestSuite) TestPing() {
	s.mock.ExpectPing()
	err := s.storage.Ping(s.ctx)
	s.Require().NoError(err)
}
