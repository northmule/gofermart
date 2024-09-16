package repository

import (
	"context"
	"github.com/northmule/gophermart/internal/app/storage"
)

type Manager struct {
	user      *UserRepository
	accrual   *AccrualRepository
	order     *OrderRepository
	withdrawn *WithdrawnRepository
	balance   *BalanceRepository
	job       *JobRepository
}

func NewManager(db storage.DBQuery, ctx context.Context) Repository {
	instance := &Manager{
		user:      NewUserRepository(db, ctx),
		accrual:   NewAccrualRepository(db, ctx),
		order:     NewOrderRepository(db, ctx),
		job:       NewJobRepository(db, ctx),
		balance:   NewBalanceRepository(db, ctx),
		withdrawn: NewWithdrawnRepository(db, ctx),
	}

	return instance
}

type Repository interface {
	User() *UserRepository
	Accrual() *AccrualRepository
	Order() *OrderRepository
	Withdrawn() *WithdrawnRepository
	Balance() *BalanceRepository
	Job() *JobRepository
}

func (m *Manager) User() *UserRepository {
	return m.user
}
func (m *Manager) Accrual() *AccrualRepository {
	return m.accrual
}
func (m *Manager) Order() *OrderRepository {
	return m.order
}
func (m *Manager) Withdrawn() *WithdrawnRepository {
	return m.withdrawn
}
func (m *Manager) Balance() *BalanceRepository {
	return m.balance
}
func (m *Manager) Job() *JobRepository {
	return m.job
}
