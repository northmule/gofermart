package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	DB         *sql.DB
	mock       sqlmock.Sqlmock
	repository *UserRepository
}

func (u *UserRepositoryTestSuite) SetupTest() {
	var err error
	u.DB, u.mock, err = sqlmock.New()

	u.mock.ExpectPrepare("select")
	u.mock.ExpectPrepare("insert")
	u.mock.ExpectPrepare("select")
	u.repository = NewUserRepository(u.DB)
	require.NoError(u.T(), err)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
