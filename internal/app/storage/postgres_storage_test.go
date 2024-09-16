package storage

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"github.com/stretchr/testify/suite"
	"testing"
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

	instance, _ := NewPostgresStorage(s.dsn, s.ctx)
	instance.DB = s.db
	s.storage = instance
}

func (s *PostgresStorageTestSuite) TearDownTest() {
	s.db.Close()
}

func (s *PostgresStorageTestSuite) TestPing() {
	s.mock.ExpectPing()
	err := s.storage.Ping()
	s.Require().NoError(err)
}

func (s *PostgresStorageTestSuite) TestTxOpen() {
	s.mock.ExpectBegin()
	err := s.storage.TxOpen()
	s.Require().NoError(err)
}

func (s *PostgresStorageTestSuite) TestTxQueryRowContext() {
	query := "SELECT"
	rows := sqlmock.NewRows([]string{})
	s.mock.ExpectBegin()
	s.storage.TxOpen()
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	_, err := s.storage.TxQueryRowContext(query)
	s.Require().NoError(err)
}

func (s *PostgresStorageTestSuite) TestTxRollback() {
	s.mock.ExpectBegin()
	err := s.storage.TxOpen()
	s.Require().NoError(err)
	s.mock.ExpectRollback()
	err = s.storage.TxRollback()
	s.Require().NoError(err)
}

func (s *PostgresStorageTestSuite) TestTxCommit() {
	s.mock.ExpectBegin()
	err := s.storage.TxOpen()
	s.Require().NoError(err)
	s.mock.ExpectCommit()
	err = s.storage.TxCommit()
	s.Require().NoError(err)
}

func TestPostgresStorageTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresStorageTestSuite))
}
