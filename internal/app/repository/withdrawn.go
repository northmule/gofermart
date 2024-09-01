package repository

import "github.com/northmule/gofermart/internal/app/storage"

type WithdrawnRepository struct {
	store storage.DBQuery
}

func NewWithdrawnRepository(store storage.DBQuery) *WithdrawnRepository {
	instance := WithdrawnRepository{
		store: store,
	}

	return &instance
}
