package repository

import "github.com/northmule/gofermart/internal/app/storage"

type Manager struct {
	User      *UserRepository
	Accrual   *AccrualRepository
	Order     *OrderRepository
	Withdrawn *WithdrawnRepository
	Balance   *BalanceRepository
}

func NewManager(db storage.DBQuery) *Manager {
	instance := &Manager{
		User:    NewUserRepository(db),
		Accrual: NewAccrualRepository(db),
		Order:   NewOrderRepository(db),
	}
	instance.Balance = NewBalanceRepository(db)
	instance.Withdrawn = NewWithdrawnRepository(db)
	return instance
}
