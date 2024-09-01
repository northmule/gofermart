package repository

import "github.com/northmule/gofermart/internal/app/storage"

type AccrualRepository struct {
	store storage.DBQuery
}

func NewAccrualRepository(store storage.DBQuery) *AccrualRepository {
	instance := AccrualRepository{
		store: store,
	}

	return &instance
}
