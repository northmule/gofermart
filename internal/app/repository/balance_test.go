package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type BalanceRepositoryTestSuite struct {
	suite.Suite
	DB         *sql.DB
	mock       sqlmock.Sqlmock
	repository *BalanceRepository
}

func (s *BalanceRepositoryTestSuite) SetupTest() {
	var err error
	s.DB, s.mock, err = sqlmock.New()

	s.mock.ExpectPrepare("select")
	s.mock.ExpectPrepare("insert")
	s.repository = NewBalanceRepository(s.DB)
	require.NoError(s.T(), err)
}
func TestBalanceRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceRepositoryTestSuite))
}

func (s *BalanceRepositoryTestSuite) TestFindOneByUserUUID() {
	userUUID := "uuid123"

	s.mock.ExpectQuery("select").
		WithArgs(userUUID).WillReturnRows(sqlmock.NewRows([]string{"b.id", "b.value", "b.updated_at", "u.id", "u.name", "u.login", " u.password", "u.created_at", "u.uuid"}).
		AddRow("1", "10", time.Now(), "2", "name", "login", "pwd", time.Now(), "uuid"))
	balance, err := s.repository.FindOneByUserUUID(context.Background(), userUUID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, balance.ID)
	require.Equal(s.T(), float64(10), balance.Value)
	require.Equal(s.T(), "uuid", balance.User.UUID)

	require.NoError(s.T(), err)
}

func (s *BalanceRepositoryTestSuite) TestCreateBalanceByUserUUID() {
	userUUID := "uuid123"

	s.mock.ExpectQuery("insert into").
		WithArgs(userUUID).WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow("1"))
	id, err := s.repository.CreateBalanceByUserUUID(context.Background(), userUUID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), int64(1), id)
}
