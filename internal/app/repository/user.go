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

type UserRepository struct {
	store          storage.DBQuery
	sqlFindByLogin *sql.Stmt
	sqlCreateUser  *sql.Stmt
	sqlFindByUUID  *sql.Stmt
}

func NewUserRepository(store storage.DBQuery) *UserRepository {
	instance := UserRepository{
		store: store,
	}
	var err error
	instance.sqlFindByLogin, err = store.Prepare(`select id, name, login, password, created_at, uuid from users where login = $1 limit 1`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}

	instance.sqlCreateUser, err = store.Prepare(`insert into users (name, login, password, uuid) values ($1, $2, $3, $4) returning id`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}

	instance.sqlFindByUUID, err = store.Prepare(`select id, name, login, password, created_at, uuid from users where uuid = $1 limit 1`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}

	return &instance
}

func (r *UserRepository) FindOneByLogin(login string) (*models.User, error) {
	user := models.User{}

	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := r.sqlFindByLogin.QueryContext(ctx, login)
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindOneByLogin(%s) произошла ошибка %s", login, err)
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindOneByLogin(%s) произошла ошибка %s", login, err)
		return nil, err
	}

	if rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt, &user.UUID)
		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindOneByLogin(%s) произошла ошибка %s", login, err)
			return nil, err
		}
	}

	return &user, nil
}

func (r *UserRepository) FindOneByUUID(uuid string) (*models.User, error) {
	user := models.User{}

	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := r.sqlFindByUUID.QueryContext(ctx, uuid)
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindOneByUUID(%s) произошла ошибка %s", uuid, err)
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindOneByUUID(%s) произошла ошибка %s", uuid, err)
		return nil, err
	}

	if rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt, &user.UUID)
		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindOneByUUID(%s) произошла ошибка %s", uuid, err)
			return nil, err
		}
	}

	return &user, nil
}

func (r *UserRepository) CreateNewUser(user models.User) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows := r.sqlCreateUser.QueryRowContext(ctx, user.Name, user.Login, user.Password, user.UUID)
	err := rows.Err()
	if err != nil {
		logger.LogSugar.Error(err)
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
