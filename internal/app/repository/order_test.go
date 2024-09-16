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

type OrderRepositoryTestSuite struct {
	suite.Suite
	DB         *sql.DB
	mock       sqlmock.Sqlmock
	repository *OrderRepository
}

func (o *OrderRepositoryTestSuite) SetupTest() {
	var err error
	o.DB, o.mock, err = sqlmock.New()

	o.mock.ExpectPrepare("select")
	o.mock.ExpectPrepare("insert")
	o.mock.ExpectPrepare("insert")
	o.mock.ExpectPrepare("select")
	o.repository = NewOrderRepository(o.DB, context.Background())
	require.NoError(o.T(), err)
}

func TestOrderRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (o *OrderRepositoryTestSuite) FindOneByNumber() {
	orderNumber := "1234"
	o.mock.ExpectQuery("select").WithArgs(orderNumber).
		WillReturnRows(sqlmock.NewRows([]string{"o.id", "o.number", "o.status", "o.created_at", "o.deleted_at", "u.id", "u.name", "u.login", "u.password", "u.created_at", "u.uuid"}).
			AddRow("12", "1002", "OK", time.Now(), time.Now(), time.Now(), "14", "userName", "login", "pwd", time.Now(), "uuid"))

	order, err := o.repository.FindOneByNumber(orderNumber)
	require.NoError(o.T(), err)
	require.Equal(o.T(), "1002", order.Number)
	require.Equal(o.T(), "OK", order.Status)
}
