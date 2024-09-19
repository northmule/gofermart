package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/northmule/gophermart/config"
	"github.com/northmule/gophermart/internal/app/repository/models"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"github.com/northmule/gophermart/internal/app/storage"
	"time"
)

type WithdrawnRepository struct {
	store                         storage.DBQuery
	sqlFindSumWithdrawnByUserUUID *sql.Stmt
	sqlFindOneByOrderID           *sql.Stmt
	sqlFindWithdrawsByUserUUID    *sql.Stmt
}

func NewWithdrawnRepository(store storage.DBQuery) *WithdrawnRepository {
	instance := WithdrawnRepository{
		store: store,
	}
	var err error
	instance.sqlFindSumWithdrawnByUserUUID, err = store.Prepare(`
																	select 
																	sum(w.value) as withdrawn
																	from withdrawals w 
																	join users u on u.id = w.user_id 
																	where u."uuid" = $1
				`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}

	instance.sqlFindOneByOrderID, err = store.Prepare(`select id, user_id, value, order_id, created_at from withdrawals where order_id = $1`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}

	instance.sqlFindWithdrawsByUserUUID, err = store.Prepare(`
												select w.id, 
												       w.user_id, 
												       w.value, 
												       w.order_id, 
												       w.created_at,
												       o.id as order_id,
												       o.number,
												       o.status,
												       o.created_at												       
												from withdrawals w
												join orders o on o.id = w.order_id
												where user_id = (select u.id from users u where u.uuid=$1 limit 1)
												order by w.id desc
												`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}

	return &instance
}

// Withdraw списание с обновлением баланса пользователя
func (wr *WithdrawnRepository) Withdraw(ctx context.Context, userUUID string, withdraw float64, orderID int) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	tx, err := wr.store.Begin()
	if err != nil {
		logger.LogSugar.Error(err)
		return 0, err
	}
	rows := tx.QueryRowContext(ctx, `insert into withdrawals (user_id, value, order_id) values ((select id from users where uuid= $1 limit 1), $2, $3) RETURNING id;`, userUUID, withdraw, orderID)
	err = rows.Err()
	if err != nil {
		err = errors.Join(err, tx.Rollback())
		logger.LogSugar.Error(err)
		return 0, err
	}

	var withdrawID int64
	err = rows.Scan(&withdrawID)
	if err != nil {
		err = errors.Join(err, tx.Rollback())
		logger.LogSugar.Error(err)
		return 0, err
	}

	rows = tx.QueryRowContext(ctx, `update user_balance set value = (value - $1) where user_id = (select id from users where uuid = $2 limit 1)`, withdraw, userUUID)
	if rows.Err() != nil {
		err = errors.Join(rows.Err(), tx.Rollback())
		logger.LogSugar.Errorf("При вызове UpdateByUserID(%d) произошла ошибка %s", userUUID, err)
		return 0, rows.Err()
	}

	if tx.Commit() != nil {
		logger.LogSugar.Error(err)
		return 0, err
	}

	return withdrawID, nil
}

func (wr *WithdrawnRepository) FindOneByOrderID(ctx context.Context, orderID int) (*models.Withdrawn, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := wr.sqlFindOneByOrderID.QueryContext(ctx, orderID)
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindOneByOrderID(%d) произошла ошибка %s", orderID, err)
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindOneByOrderID(%d) произошла ошибка %s", orderID, err)
		return nil, err
	}
	var withdraw models.Withdrawn
	if rows.Next() {
		err := rows.Scan(&withdraw.ID, &withdraw.UserID, &withdraw.Value, &withdraw.OrderID, &withdraw.CreatedAt)
		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindOneByOrderID(%d) произошла ошибка %s", orderID, err)
			return nil, err
		}
	}
	return &withdraw, nil
}

func (wr *WithdrawnRepository) FindSumWithdrawnByUserUUID(ctx context.Context, userUUID string) (float64, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := wr.sqlFindSumWithdrawnByUserUUID.QueryContext(ctx, userUUID)
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindOneByUserUUID(%s) произошла ошибка %s", userUUID, err)
		return 0, err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindOneByUserUUID(%s) произошла ошибка %s", userUUID, err)
		return 0, err
	}
	var sum sql.NullFloat64
	if rows.Next() {
		err := rows.Scan(&sum)
		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindOneByUserUUID(%s) произошла ошибка %s", userUUID, err)
			return 0, err
		}
	}
	return sum.Float64, nil
}

func (wr *WithdrawnRepository) FindWithdrawsByUserUUID(ctx context.Context, userUUID string) ([]models.Withdrawn, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := wr.sqlFindWithdrawsByUserUUID.QueryContext(ctx, userUUID)
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindWithdrawsByUserUUID(%s) произошла ошибка %s", userUUID, err)
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindWithdrawsByUserUUID(%s) произошла ошибка %s", userUUID, err)
		return nil, err
	}
	var withdraws []models.Withdrawn
	for rows.Next() {
		order := models.Order{}
		withdraw := models.Withdrawn{}
		err := rows.Scan(&withdraw.ID, &withdraw.UserID, &withdraw.Value, &withdraw.OrderID, &withdraw.CreatedAt, &order.ID, &order.Number, &order.Status, &order.CreatedAt)

		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindWithdrawsByUserUUID(%s) произошла ошибка %s", userUUID, err)
			return nil, err
		}
		withdraw.Order = &order
		withdraws = append(withdraws, withdraw)
	}

	return withdraws, nil
}
