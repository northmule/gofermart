package repository

import (
	"context"
	"database/sql"
	"github.com/northmule/gofermart/config"
	"github.com/northmule/gofermart/internal/app/repository/models"
	"github.com/northmule/gofermart/internal/app/services/logger"
	"github.com/northmule/gofermart/internal/app/storage"
	"time"
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
	instance.sqlFindByUserUUID, err = store.Prepare(`
				select b.id, b.value, b.updated_at, u.id, u.name, u.login, u.password, u.created_at, u.uuid from user_balance as b
				join users as u on u.id = b.user_id
				where u.uuid = $1
				limit 1
				`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}

	return &instance
}

func (br *BalanceRepository) FindOneByUserUUID(userUUID string) (*models.Balance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := br.sqlFindByUserUUID.QueryContext(ctx, userUUID)
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindOneByUserUUID(%s) произошла ошибка %s", userUUID, err)
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindOneByUserUUID(%s) произошла ошибка %s", userUUID, err)
		return nil, err
	}
	var balance models.Balance
	var user models.User
	if rows.Next() {
		err := rows.Scan(&balance.ID, &balance.Value, &balance.UpdatedAt, &user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt, &user.UUID)
		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindOneByUserUUID(%s) произошла ошибка %s", userUUID, err)
			return nil, err
		}
	}
	balance.User = user
	return &balance, nil
}
