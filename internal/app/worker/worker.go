package worker

import (
	"context"
	"errors"
	"github.com/northmule/gophermart/internal/accrual/api/client"
	"github.com/northmule/gophermart/internal/app/constants"
	"github.com/northmule/gophermart/internal/app/repository"
	"github.com/northmule/gophermart/internal/app/repository/models"
	"github.com/northmule/gophermart/internal/app/services/logger"
)

const maxNumberAttempts = 5

type Worker struct {
	manager        repository.Repository
	accrualService AccrualClientInterface
	jobChan        chan job
	workerNum      int
}

func NewWorker(manager repository.Repository, accrualService AccrualClientInterface, workerNum int) *Worker {
	instance := Worker{
		manager:        manager,
		accrualService: accrualService,
		jobChan:        make(chan job, 1),
		workerNum:      workerNum,
	}

	return &instance
}

type AccrualClientInterface interface {
	SendOrderNumber(ctx context.Context, orderNumber string) (*client.ResponseAccrual, error)
}

type job struct {
	jobRun models.Job
}

func (w *Worker) Run(ctx context.Context) {
	for i := 1; i <= w.workerNum; i++ {
		go w.worker(ctx, i, w.jobChan)
	}

	go w.producer(ctx)
}

func (w *Worker) producer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.LogSugar.Info("Останавливаю producer по сигналу")
			return
		default:
			jobsForRun, err := w.manager.Job().GetJobForRun(ctx)
			if err != nil {
				logger.LogSugar.Error(err.Error())
				break
			}
			for _, jobForRun := range jobsForRun {
				logger.LogSugar.Infof("В работу взята задча на запрос к внешнему сервису с номером заказа %s", jobForRun.OrderNumber)
				w.jobChan <- job{jobRun: jobForRun}
			}
		}
	}
}

func (w *Worker) worker(ctx context.Context, num int, jobCh <-chan job) {
	var errorNoContent *client.ErrorNoContent
	var errorTooManyRequests *client.ErrorTooManyRequests
	var errorInternalServerError *client.ErrorInternalServerError
	var errorUndefined *client.ErrorUndefined
	logger.LogSugar.Infof("Запуск worker %d из %d", num, w.workerNum)
	for {
		select {
		case <-ctx.Done():
			logger.LogSugar.Infof("Останавливаю worker %d из %d по сигналу", num, w.workerNum)
			return
		case item := <-jobCh:
			response, err := w.accrualService.SendOrderNumber(ctx, item.jobRun.OrderNumber)

			// Условно ошибки вызванные отсутсвием данных. Запрашиваем maxNumberAttempts раз, потом удаляем задачу
			if errors.As(err, &errorUndefined) || errors.As(err, &errorInternalServerError) || errors.As(err, &errorNoContent) || err != nil {
				logger.LogSugar.Info(err)
				if item.jobRun.RunCnt > maxNumberAttempts {
					err = w.manager.Job().DeleteJobByOrderNumber(ctx, item.jobRun.OrderNumber)
					logger.LogSugar.Infof("Удалил задачу по обработке заказа %s", item.jobRun.OrderNumber)
				} else {
					err = w.manager.Job().UpdateJobByOrderNumber(ctx, item.jobRun.OrderNumber)
				}
				if err != nil {
					logger.LogSugar.Info(err)
				}
				break
			}

			// Условно положительные ошиби и статусы, откладываем выполнение пока не будет положительныйо ответ
			if errors.As(err, &errorTooManyRequests) ||
				response.Status == constants.OrderStatusNew ||
				response.Status == constants.OrderStatusProcessing {
				if err = w.manager.Job().UpdateJobByOrderNumber(ctx, item.jobRun.OrderNumber); err != nil {
					logger.LogSugar.Info(err)
				}
				break
			}

			// Положительный ответ. Обновляем начисления, статус заказа, удаляем задачу
			if response.Status == constants.OrderStatusInvalid || response.Status == constants.OrderStatusProcessed {
				logger.LogSugar.Infof("Обновляю информацию о начислениях по заказу %s. Будет начисленно %f", response.Order, response.Accrual)
				if err = w.manager.Accrual().UpdateTxByOrderNumber(ctx, response.Order, response.Status, response.Accrual); err != nil {
					logger.LogSugar.Error(err)
				}
				if err = w.manager.Job().DeleteJobByOrderNumber(ctx, item.jobRun.OrderNumber); err != nil {
					logger.LogSugar.Error(err)
				}
				logger.LogSugar.Infof("Удалил задачу по обработке заказа %s", item.jobRun.OrderNumber)
			}
		}
	}
}
