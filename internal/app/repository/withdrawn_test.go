package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type WithdrawnRepositoryTestSuite struct {
	suite.Suite
	DB         *sql.DB
	mock       sqlmock.Sqlmock
	repository *WithdrawnRepository
}

func (w *WithdrawnRepositoryTestSuite) SetupTest() {
	var err error
	w.DB, w.mock, err = sqlmock.New()

	w.mock.ExpectPrepare("select")
	w.mock.ExpectPrepare("select")
	w.mock.ExpectPrepare("select")
	w.repository = NewWithdrawnRepository(w.DB, context.Background())
	require.NoError(w.T(), err)
}

func TestWithdrawnRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(WithdrawnRepositoryTestSuite))
}
