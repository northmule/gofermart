package worker

import (
	"github.com/northmule/gofermart/internal/app/services/logger"
)

const workerNum = 3

type Worker struct {
	deleter Deleter
	jobChan chan job
}

func NewWorker(deleter Deleter) *Worker {
	instance := Worker{
		deleter: deleter,
	}

	instance.jobChan = make(chan job, 5)

	for i := 0; i < workerNum; i++ {
		go instance.worker(instance.jobChan)
	}

	return &instance
}

type Deleter interface {
	SoftDeletedShortURL(userUUID string, shortURL ...string) error
}

type job struct {
	uuid string
	url  string
}

func (w *Worker) Del(userUUID string, input []string) {

	go w.producer(userUUID, input, w.jobChan)
}

func (w *Worker) producer(uuid string, urls []string, jobCh chan<- job) {
	for _, url := range urls {
		jobCh <- job{uuid: uuid, url: url}
	}
}

func (w *Worker) worker(jobCh <-chan job) {
	for item := range jobCh {
		err := w.deleter.SoftDeletedShortURL(item.uuid, item.url)
		if err != nil {
			logger.LogSugar.Infof(err.Error())
		}
	}
}
