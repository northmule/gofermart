package accrual

import (
	"context"
	"github.com/northmule/gophermart/internal/app/repository"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"go.uber.org/zap"
)

type AccrualService struct {
	manager repository.Repository
}

func NewAccrualService(manager repository.Repository) *AccrualService {
	instance := &AccrualService{
		manager: manager,
	}
	return instance
}

func (a *AccrualService) CreateJobToProcessNewOrder(ctx context.Context, orderNumber string, userUUID string) error {

	_, err := a.manager.Accrual().CreateAccrualByOrderNumberAndUserUUID(ctx, orderNumber, userUUID)
	logger.LogSugar.Info("Создаю информацию о нулевом списании по заказу", zap.String("number", orderNumber))
	if err != nil {
		return err
	}

	return nil
}
