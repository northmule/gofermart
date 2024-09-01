package repository

import "github.com/northmule/gofermart/internal/app/storage"

type UserRepository struct {
	store storage.DBQuery
}

func NewUserRepository(store storage.DBQuery) *UserRepository {
	instance := UserRepository{
		store: store,
	}

	return &instance
}
