package repository

<<<<<<< HEAD
import (
	"context"
	"database/sql"
	"github.com/northmule/gofermart/config"
	"github.com/northmule/gofermart/internal/app/repository/models"
	"github.com/northmule/gofermart/internal/app/services/logger"
	"github.com/northmule/gofermart/internal/app/storage"
	"time"
)

type OrderRepository struct {
	store                   storage.DBQuery
	sqlFindByNumber         *sql.Stmt
	sqlCreateOrder          *sql.Stmt
	sqlLinkOrderToUser      *sql.Stmt
	sqlFindOrdersByUserUUID *sql.Stmt
=======
import "github.com/northmule/gofermart/internal/app/storage"

type OrderRepository struct {
	store storage.DBQuery
>>>>>>> 94746e2 (базовая структура)
}

func NewOrderRepository(store storage.DBQuery) *OrderRepository {
	instance := OrderRepository{
		store: store,
	}

<<<<<<< HEAD
	var err error
	instance.sqlFindByNumber, err = store.Prepare(`select o.id, o.number, o.status, o.created_at, o.deleted_at,
       													u.id, u.name, u.login, u.password, u.created_at, u.uuid
															from orders as o
                                                  join user_orders uo on uo.order_id = o.id
                                                  join users u on u.id = uo.user_id
                                                  where number = $1 limit 1`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}

	instance.sqlCreateOrder, err = store.Prepare(`insert into orders (number, status) values ($1, $2) RETURNING id;`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}
	instance.sqlLinkOrderToUser, err = store.Prepare(`insert into user_orders (user_id, order_id) values ($1, $2) RETURNING id;`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}
	instance.sqlFindOrdersByUserUUID, err = store.Prepare(`select 
																	o.id, 
																	o."number", 
																	o.status, 
																	o.created_at, 
																	o.deleted_at,
																	sum(a.value) as accrual
																	from orders o 
																	join user_orders uo on uo.order_id = o.id 
																	join users u on u.id = uo.user_id 
																	left join accruals a on a.order_id = o.id 
																	where u."uuid" = $1
																	group by o.id
																	order by o.id desc`,
	)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}

	return &instance
}

func (o *OrderRepository) FindOneByNumber(number string) (*models.Order, error) {
	order := models.Order{}
	user := models.User{}

	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := o.sqlFindByNumber.QueryContext(ctx, number)
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindOneByNumber(%s) произошла ошибка %s", number, err)
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindOneByNumber(%s) произошла ошибка %s", number, err)
		return nil, err
	}

	if rows.Next() {
		err := rows.Scan(&order.ID, &order.Number, &order.Status, &order.CreatedAt, &order.DeletedAt, &user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt, &user.UUID)
		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindOneByNumber(%s) произошла ошибка %s", number, err)
			return nil, err
		}
	}
	order.User = user

	return &order, nil
}

func (o *OrderRepository) Save(order models.Order, userID int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows := o.sqlCreateOrder.QueryRowContext(ctx, order.Number, order.Status)
	err := rows.Err()
	if err != nil {
		logger.LogSugar.Error(err)
		return 0, err
	}

	var orderID int64
	err = rows.Scan(&orderID)
	if err != nil {
		logger.LogSugar.Error(err)
		return 0, err
	}
	rows = o.sqlLinkOrderToUser.QueryRowContext(ctx, userID, orderID)
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Error(err)
		return 0, err
	}

	return orderID, nil
}

func (o *OrderRepository) FindOrdersByUserUUID(userUUID string) (*[]models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := o.sqlFindOrdersByUserUUID.QueryContext(ctx, userUUID)
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindOrdersByUserUUID(%s) произошла ошибка %s", userUUID, err)
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindOrdersByUserUUID(%s) произошла ошибка %s", userUUID, err)
		return nil, err
	}
	var orders []models.Order
	for rows.Next() {
		order := models.Order{}
		err := rows.Scan(&order.ID, &order.Number, &order.Status, &order.CreatedAt, &order.DeletedAt, &order.Accrual)
		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindOrdersByUserUUID(%s) произошла ошибка %s", userUUID, err)
			return nil, err
		}
		orders = append(orders, order)
	}

	return &orders, nil
}
=======
	return &instance
}
>>>>>>> 94746e2 (базовая структура)
