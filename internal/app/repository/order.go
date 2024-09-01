package repository

import "github.com/northmule/gofermart/internal/app/storage"

type OrderRepository struct {
	store storage.DBQuery
}

func NewOrderRepository(store storage.DBQuery) *OrderRepository {
	instance := OrderRepository{
		store: store,
	}

	return &instance
}
