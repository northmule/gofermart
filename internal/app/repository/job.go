package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/northmule/gophermart/config"
	"github.com/northmule/gophermart/internal/app/repository/models"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"github.com/northmule/gophermart/internal/app/storage"
	"time"
)

type JobRepository struct {
	store                     storage.DBQuery
	sqlGetJobForRun           *sql.Stmt
	sqlCreateJobByOrderNumber *sql.Stmt
	sqlUpdateJobByOrderNumber *sql.Stmt
	sqlDeleteJobByOrderNumber *sql.Stmt
}

func NewJobRepository(store storage.DBQuery) *JobRepository {
	instance := JobRepository{
		store: store,
	}
	var err error

	instance.sqlGetJobForRun, err = store.Prepare(`select id, order_number, created_at, updated_at, next_run, run_cnt from jobs_order where next_run < now() and is_work = false order by id asc`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}

	instance.sqlCreateJobByOrderNumber, err = store.Prepare(`insert into jobs_order (order_number) values ($1) returning id;`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}

	instance.sqlUpdateJobByOrderNumber, err = store.Prepare(`update jobs_order set next_run = now() + interval '10 sec', run_cnt = run_cnt + 1, is_work = false where order_number = $1`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}

	instance.sqlDeleteJobByOrderNumber, err = store.Prepare(`delete from jobs_order where order_number = $1`)
	if err != nil {
		logger.LogSugar.Error(err)
		return nil
	}

	return &instance
}

func (jr *JobRepository) GetJobForRun(ctx context.Context) (*[]models.Job, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()

	tx, err := jr.store.Begin()
	if err != nil {
		logger.LogSugar.Error(err)
		return nil, err
	}
	rows, err := tx.QueryContext(ctx, `select id, order_number, created_at, updated_at, next_run, run_cnt from jobs_order where next_run < now() and is_work = false order by id asc`)
	if err != nil {
		err = errors.Join(err, tx.Rollback())
		logger.LogSugar.Errorf("При вызове GetJobForRun() произошла ошибка %s", err)
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		err = errors.Join(err, rows.Err(), tx.Rollback())
		logger.LogSugar.Errorf("При вызове GetJobForRun произошла ошибка %s", err)
		return nil, err
	}
	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(&job.ID, &job.OrderNumber, &job.CreatedAt, &job.UpdatedAt, &job.NextRun, &job.RunCnt)
		if err != nil {
			err = errors.Join(err, rows.Err(), tx.Rollback())
			logger.LogSugar.Errorf("При обработке значений в GetJobForRun произошла ошибка %s", err)
			return nil, err
		}
		jobs = append(jobs, job)

	}

	for _, job := range jobs {
		rows, err := tx.QueryContext(ctx, `update jobs_order set is_work = true where id = $1`, job.ID)
		if rows.Err() != nil {
			err = errors.Join(rows.Err(), tx.Rollback())
			logger.LogSugar.Errorf("При вызове GetJobForRun произошла ошибка %s", err)
			return nil, err
		}
		if err != nil {
			err = errors.Join(err, rows.Err(), tx.Rollback())
			logger.LogSugar.Errorf("При вызове GetJobForRun произошла ошибка %s", err)
			return nil, err
		}
	}

	if tx.Commit() != nil {
		logger.LogSugar.Error(err)
		return nil, err
	}

	return &jobs, nil
}

func (jr *JobRepository) CreateJobByOrderNumber(ctx context.Context, orderNumber string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows := jr.sqlCreateJobByOrderNumber.QueryRowContext(ctx, orderNumber)

	err := rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове CreateJobByOrderNumber(%s) произошла ошибка %s", orderNumber, err)
		return 0, err
	}

	var id int64
	err = rows.Scan(&id)
	if err != nil {
		logger.LogSugar.Error(err)
		return 0, err
	}

	return id, nil
}

func (jr *JobRepository) UpdateJobByOrderNumber(ctx context.Context, orderNumber string) error {
	ctx, cancel := context.WithTimeout(ctx, config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := jr.sqlUpdateJobByOrderNumber.QueryContext(ctx, orderNumber)
	if err != nil {
		logger.LogSugar.Errorf("При вызове UpdateJobByOrderNumber(%s) произошла ошибка %s", orderNumber, err)
		return err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове UpdateJobByOrderNumber(%s) произошла ошибка %s", orderNumber, err)
		return err
	}

	return nil
}

func (jr *JobRepository) DeleteJobByOrderNumber(ctx context.Context, orderNumber string) error {
	ctx, cancel := context.WithTimeout(ctx, config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := jr.sqlDeleteJobByOrderNumber.QueryContext(ctx, orderNumber)
	if err != nil {
		logger.LogSugar.Errorf("При вызове DeleteJobByOrderNumber(%s) произошла ошибка %s", orderNumber, err)
		return err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове DeleteJobByOrderNumber(%s) произошла ошибка %s", orderNumber, err)
		return err
	}

	return nil
}
