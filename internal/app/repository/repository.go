package repository

import "github.com/northmule/gophermart/internal/app/storage"

type Manager struct {
	User      *UserRepository
	Accrual   *AccrualRepository
	Order     *OrderRepository
	Withdrawn *WithdrawnRepository
	Balance   *BalanceRepository
	Job       *JobRepository
}

func NewManager(db storage.DBQuery) *Manager {
	instance := &Manager{
		User:    NewUserRepository(db),
		Accrual: NewAccrualRepository(db),
		Order:   NewOrderRepository(db),
		Job:     NewJobRepository(db),
	}
	instance.Balance = NewBalanceRepository(db)
	instance.Withdrawn = NewWithdrawnRepository(db)
	return instance
}
