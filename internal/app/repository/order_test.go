package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/northmule/gophermart/internal/app/repository/models"
	"github.com/shopspring/decimal"
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
	o.mock.ExpectPrepare("select")
	o.repository = NewOrderRepository(o.DB)
	require.NoError(o.T(), err)
}

func TestOrderRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (o *OrderRepositoryTestSuite) TestFindOneByNumber() {
	orderNumber := "1234"
	o.mock.ExpectQuery("select").WithArgs(orderNumber).
		WillReturnRows(sqlmock.NewRows([]string{"o.id", "o.number", "o.status", "o.created_at", "o.deleted_at", "u.id", "u.login", "u.password", "u.created_at", "u.uuid"}).
			AddRow("12", "1002", "OK", time.Now(), time.Now(), "14", "login", "pwd", time.Now(), "uuid"))

	order, err := o.repository.FindOneByNumber(context.Background(), orderNumber)
	require.NoError(o.T(), err)
	require.Equal(o.T(), "1002", order.Number)
	require.Equal(o.T(), "OK", order.Status)
}

func (o *OrderRepositoryTestSuite) TestFindOrdersByUserUUID() {
	userUUID := "uuid-user"
	o.mock.ExpectQuery("select").WithArgs(userUUID).
		WillReturnRows(sqlmock.NewRows([]string{"o.id", "o.number", "o.status", "o.created_at", "o.deleted_at", "o.accrual"}).
			AddRow("12", "1002", "OK", time.Now(), nil, "134"))

	orders, err := o.repository.FindOrdersByUserUUID(context.Background(), userUUID)
	require.Len(o.T(), orders, 1)
	require.NoError(o.T(), err)
	order := orders[0]
	require.Equal(o.T(), decimal.NewFromFloat(134), order.Accrual)
	require.Equal(o.T(), "OK", order.Status)
}

func (o *OrderRepositoryTestSuite) TestSave() {
	userUUID := "uuid-user"
	order := models.Order{
		Number: "123",
		Status: "Ok",
	}
	expectedID := int64(2)
	o.mock.ExpectBegin()
	o.mock.ExpectQuery("insert into orders").WithArgs(order.Number, order.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	o.mock.ExpectQuery("insert into user_orders").WithArgs(userUUID, int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(112))
	o.mock.ExpectCommit()

	id, err := o.repository.Save(context.Background(), order, userUUID)
	require.NoError(o.T(), err)
	require.Equal(o.T(), expectedID, id)
}
