package repository

import "github.com/northmule/gofermart/internal/app/storage"

type Manager struct {
	User      *UserRepository
	Accrual   *AccrualRepository
	Order     *OrderRepository
	Withdrawn *WithdrawnRepository
<<<<<<< HEAD
	Balance   *BalanceRepository
=======
>>>>>>> 94746e2 (базовая структура)
}

func NewManager(db storage.DBQuery) *Manager {
	instance := &Manager{
<<<<<<< HEAD
		User:    NewUserRepository(db),
		Accrual: NewAccrualRepository(db),
		Order:   NewOrderRepository(db),
	}
	instance.Balance = NewBalanceRepository(db)
	instance.Withdrawn = NewWithdrawnRepository(db)
=======
		User:      NewUserRepository(db),
		Accrual:   NewAccrualRepository(db),
		Order:     NewOrderRepository(db),
		Withdrawn: NewWithdrawnRepository(db),
	}

>>>>>>> 94746e2 (базовая структура)
	return instance
}
