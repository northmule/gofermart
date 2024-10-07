package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AccrualRepositoryTestSuite struct {
	suite.Suite
	DB   *sql.DB
	mock sqlmock.Sqlmock
}

func (suite *AccrualRepositoryTestSuite) SetupTest() {
	var err error
	suite.DB, suite.mock, err = sqlmock.New()
	require.NoError(suite.T(), err)
}

func (suite *AccrualRepositoryTestSuite) TestCreateAccrualByOrderNumberAndUserUUID() {
	suite.mock.ExpectPrepare("insert")
	ar := NewAccrualRepository(suite.DB)

	orderNumber := "12345"
	userUUID := "uuid123"
	expectedID := int64(1)

	suite.mock.ExpectQuery("insert into").
		WithArgs(orderNumber, userUUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	id, err := ar.CreateAccrualByOrderNumberAndUserUUID(context.Background(), orderNumber, userUUID)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), expectedID, id)

	err = suite.mock.ExpectationsWereMet()
	require.NoError(suite.T(), err)
}

func (suite *AccrualRepositoryTestSuite) TestUpdateTxByOrderNumber() {
	suite.mock.ExpectPrepare("insert")
	ar := NewAccrualRepository(suite.DB)

	orderNumber := "12345"
	orderStatus := "ok"
	accrual := decimal.NewFromFloat(50)

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery("update accruals").
		WithArgs(orderNumber, accrual).
		WillReturnRows(sqlmock.NewRows([]string{""}))

	suite.mock.ExpectQuery("update orders").
		WithArgs(orderStatus, orderNumber).
		WillReturnRows(sqlmock.NewRows([]string{""}))

	suite.mock.ExpectQuery("update user_balance").
		WithArgs(accrual, orderNumber).
		WillReturnRows(sqlmock.NewRows([]string{""}))

	suite.mock.ExpectCommit()

	err := ar.UpdateTxByOrderNumber(context.Background(), orderNumber, orderStatus, accrual)
	require.NoError(suite.T(), err)

	err = suite.mock.ExpectationsWereMet()
	require.NoError(suite.T(), err)
}

func TestAccrualRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AccrualRepositoryTestSuite))
}
