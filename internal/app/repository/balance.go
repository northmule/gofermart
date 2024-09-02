package repository

import (
	"database/sql"
	"github.com/northmule/gofermart/internal/app/services/logger"
	"github.com/northmule/gofermart/internal/app/storage"
)

type BalanceRepository struct {
	store             storage.DBQuery
	sqlFindByUserUUID *sql.Stmt
}

func NewBalanceRepository(store storage.DBQuery) *BalanceRepository {
	instance := BalanceRepository{
		store: store,
	}

	var err error
	instance.sqlFindByUserUUID, err = store.Prepare(`select o.id, o.number, o.status, o.created_at, o.deleted_at,
       													u.id, u.name, u.login, u.password, u.created_at, u.uuid
															from orders as o
                                                  join user_orders uo on uo.order_id = o.id
                                                  join users u on u.id = uo.user_id
                                                  where number = $1 limit 1`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}

	return &instance
}
