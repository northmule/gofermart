package repository

import (
	"context"
	"github.com/northmule/gophermart/internal/app/storage"
)

type Manager struct {
	User      *UserRepository
	Accrual   *AccrualRepository
	Order     *OrderRepository
	Withdrawn *WithdrawnRepository
	Balance   *BalanceRepository
	Job       *JobRepository
}

func NewManager(db storage.DBQuery, ctx context.Context) *Manager {
	instance := &Manager{
		User:      NewUserRepository(db, ctx),
		Accrual:   NewAccrualRepository(db, ctx),
		Order:     NewOrderRepository(db, ctx),
		Job:       NewJobRepository(db, ctx),
		Balance:   NewBalanceRepository(db, ctx),
		Withdrawn: NewWithdrawnRepository(db, ctx),
	}

	return instance
}
