package job

import (
	"context"
	"errors"
	"fmt"
	"github.com/northmule/gophermart/internal/app/repository"
	"github.com/northmule/gophermart/internal/app/services/logger"
)

type JobService struct {
	manager repository.Repository
}

func NewJobService(manager repository.Repository) *JobService {
	instance := &JobService{
		manager: manager,
	}
	return instance
}

func (j *JobService) CreateJobToProcessNewOrder(ctx context.Context, orderNumber string) error {

	_, err := j.manager.Job().CreateJobByOrderNumber(ctx, orderNumber)

	if err != nil {
		errors.Join(err, fmt.Errorf("ошибка создания задания на обработку заказа с номером %s", orderNumber))
		return err
	}
	logger.LogSugar.Infof("Создал задачу для запроса начисленных балов для заказа %s", orderNumber)
	return nil
}
