package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/northmule/gophermart/config"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"github.com/northmule/gophermart/internal/app/storage"
	"time"
)

type AccrualRepository struct {
	store                     storage.DBQuery
	sqlCreateAccrualZeroValue *sql.Stmt
}

func NewAccrualRepository(store storage.DBQuery) *AccrualRepository {
	instance := AccrualRepository{
		store: store,
	}
	var err error

	instance.sqlCreateAccrualZeroValue, err = store.Prepare(`insert into accruals (order_id, value, user_id) values((select o.id from orders o where o.number = $1 limit 1), 0, (select u.id from users u where u.uuid = $2 limit 1)) returning id`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}

	return &instance
}

func (ar *AccrualRepository) CreateAccrualByOrderNumberAndUserUUID(ctx context.Context, orderNumber string, userUUID string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows := ar.sqlCreateAccrualZeroValue.QueryRowContext(ctx, orderNumber, userUUID)
	err := rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове CreateAccrualByOrderNumberAndUserUUID(%s) произошла ошибка %s", userUUID, err)
		return 0, err
	}

	var id int64
	err = rows.Scan(&id)
	if err != nil {
		logger.LogSugar.Error(err)
		return 0, err
	}
	return id, nil
}

func (ar *AccrualRepository) UpdateTxByOrderNumber(ctx context.Context, orderNumber string, orderStatus string, accrual float64) error {
	ctx, cancel := context.WithTimeout(ctx, config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	tx, err := ar.store.Begin()
	if err != nil {
		logger.LogSugar.Error(err)
		return err
	}
	rows := tx.QueryRowContext(ctx, `update accruals set value = $2 where order_id = (select id from orders o where o.number = $1 limit 1)`, orderNumber, accrual)
	err = rows.Err()
	if err != nil {
		err = errors.Join(err, tx.Rollback())
		logger.LogSugar.Errorf("При вызове UpdateByOrderNumber(%s, %s, %f) произошла ошибка %s", orderNumber, orderStatus, accrual, err)
		return err
	}

	rows = tx.QueryRowContext(ctx, `update orders set status = $1 where number = $2`, orderStatus, orderNumber)
	if rows.Err() != nil {
		err = errors.Join(rows.Err(), tx.Rollback())
		logger.LogSugar.Errorf("При вызове UpdateByOrderNumber(%s, %s, %f) произошла ошибка %s", orderNumber, orderStatus, accrual, err)
		return err
	}

	rows = tx.QueryRowContext(ctx, `update user_balance set value = value + $1 where user_id = (select a.user_id from accruals a where a.order_id = (select o.id from orders o where o.number = $2 limit 1) limit 1)`, accrual, orderNumber)
	if rows.Err() != nil {
		err = errors.Join(rows.Err(), tx.Rollback())
		logger.LogSugar.Errorf("При вызове UpdateByOrderNumber(%s, %s, %f) произошла ошибка %s", orderNumber, orderStatus, accrual, err)
		return err
	}

	if tx.Commit() != nil {
		logger.LogSugar.Error(err)
		return err
	}

	return nil
}
