package repository

import "github.com/northmule/gofermart/internal/app/storage"

type Manager struct {
	User      *UserRepository
	Accrual   *AccrualRepository
	Order     *OrderRepository
	Withdrawn *WithdrawnRepository
}

func NewManager(db storage.DBQuery) *Manager {
	instance := &Manager{
		User:      NewUserRepository(db),
		Accrual:   NewAccrualRepository(db),
		Order:     NewOrderRepository(db),
		Withdrawn: NewWithdrawnRepository(db),
	}

	return instance
}
