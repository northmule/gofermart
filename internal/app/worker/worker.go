package job

import (
	"errors"
	"github.com/northmule/gofermart/internal/accrual/api/client"
	"github.com/northmule/gofermart/internal/app/constants"
	"github.com/northmule/gofermart/internal/app/repository"
	"github.com/northmule/gofermart/internal/app/repository/models"
	"github.com/northmule/gofermart/internal/app/services/logger"
)

const workerNum = 3
const maxNumberAttempts = 5

type Worker struct {
	manager        *repository.Manager
	accrualService *client.AccrualClient
	jobChan        chan job
}

func NewWorker(manager *repository.Manager, accrualService *client.AccrualClient) *Worker {
	instance := Worker{
		manager:        manager,
		accrualService: accrualService,
	}

	instance.jobChan = make(chan job, 1)

	for i := 0; i < workerNum; i++ {
		go instance.worker(instance.jobChan)
	}

	go instance.producer()

	return &instance
}

type job struct {
	jobRun models.Job
}

func (w *Worker) producer() {
	for {
		jobsForRun, err := w.manager.Job.GetJobForRun()
		if err != nil {
			logger.LogSugar.Error(err.Error())
			return
		}
		for _, jobForRun := range *jobsForRun {
			w.jobChan <- job{jobRun: jobForRun}
		}
	}
}

func (w *Worker) worker(jobCh <-chan job) {
	var errorNoContent *client.ErrorNoContent
	var errorTooManyRequests *client.ErrorTooManyRequests
	var errorInternalServerError *client.ErrorInternalServerError
	var errorUndefined *client.ErrorUndefined

	for item := range jobCh {
		response, err := w.accrualService.SendOrderNumber(item.jobRun.OrderNumber)

		if errors.As(err, &errorNoContent) {
			logger.LogSugar.Info(err.Error())
			if item.jobRun.RunCnt > maxNumberAttempts {
				err = w.manager.Job.DeleteJobByOrderNumber(item.jobRun.OrderNumber)
			} else {
				err = w.manager.Job.UpdateJobByOrderNumber(item.jobRun.OrderNumber)
			}
			if err != nil {
				logger.LogSugar.Info(err.Error())
			}
			continue
		}

		if errors.As(err, &errorTooManyRequests) {
			logger.LogSugar.Info(err.Error())
			err = w.manager.Job.UpdateJobByOrderNumber(item.jobRun.OrderNumber)
			if err != nil {
				logger.LogSugar.Info(err.Error())
			}
			continue
		}

		if errors.As(err, &errorInternalServerError) {
			logger.LogSugar.Info(err.Error())
			if item.jobRun.RunCnt > maxNumberAttempts {
				err = w.manager.Job.DeleteJobByOrderNumber(item.jobRun.OrderNumber)
			} else {
				err = w.manager.Job.UpdateJobByOrderNumber(item.jobRun.OrderNumber)
			}
			if err != nil {
				logger.LogSugar.Info(err.Error())
			}
			continue
		}

		if errors.As(err, &errorUndefined) {
			logger.LogSugar.Info(err.Error())
			if item.jobRun.RunCnt > maxNumberAttempts {
				err = w.manager.Job.DeleteJobByOrderNumber(item.jobRun.OrderNumber)
			} else {
				err = w.manager.Job.UpdateJobByOrderNumber(item.jobRun.OrderNumber)
			}
			if err != nil {
				logger.LogSugar.Info(err.Error())
			}
			continue
		}
		// Не пойманная ошибка
		if err != nil {
			err = w.manager.Job.UpdateJobByOrderNumber(item.jobRun.OrderNumber)
			err = errors.Join(err)
			logger.LogSugar.Error(err.Error())
			continue
		}
		// Откладываем задачу
		if response.Status == constants.OrderStatusNew || response.Status == constants.OrderStatusProcessing {
			err = w.manager.Job.UpdateJobByOrderNumber(item.jobRun.OrderNumber)

			if err != nil {
				logger.LogSugar.Info(err.Error())
			}
			continue
		}
		// Обновляем начисления, статус заказа, удаляем задачу
		if response.Status == constants.OrderStatusInvalid || response.Status == constants.OrderStatusProcessed {
			err = w.manager.Accrual.UpdateTxByOrderNumber(response.Order, response.Status, response.Accrual)
			if err != nil {
				logger.LogSugar.Error(err.Error())
			}
			err = w.manager.Job.DeleteJobByOrderNumber(item.jobRun.OrderNumber)
			if err != nil {
				logger.LogSugar.Error(err.Error())
			}
		}

	}
}
