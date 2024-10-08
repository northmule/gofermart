package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	return &instance
}

func (jr *JobRepository) GetJobForRun(ctx context.Context) ([]models.Job, error) {
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
		return nil, fmt.Errorf("при вызове GetJobForRun() произошла ошибка %w", err)
	}

	err = rows.Err()
	if err != nil {
		err = errors.Join(err, rows.Err(), tx.Rollback())
		return nil, fmt.Errorf("при вызове GetJobForRun произошла ошибка %w", err)
	}
	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err = rows.Scan(&job.ID, &job.OrderNumber, &job.CreatedAt, &job.UpdatedAt, &job.NextRun, &job.RunCnt)
		if err != nil {
			err = errors.Join(err, rows.Err(), tx.Rollback())
			return nil, fmt.Errorf("при обработке значений в GetJobForRun произошла ошибка %w", err)
		}

		jobs = append(jobs, job)
	}

	for _, job := range jobs {
		rows, err = tx.QueryContext(ctx, `update jobs_order set is_work = true where id = $1`, job.ID)
		if rows.Err() != nil {
			err = errors.Join(rows.Err(), tx.Rollback())
			return nil, fmt.Errorf("при вызове GetJobForRun произошла ошибка %w", err)
		}
		if err != nil {
			err = errors.Join(err, rows.Err(), tx.Rollback())
			return nil, fmt.Errorf("при вызове GetJobForRun произошла ошибка %w", err)
		}
	}

	if tx.Commit() != nil {
		logger.LogSugar.Error(err)
		return nil, err
	}

	return jobs, nil
}

func (jr *JobRepository) CreateJobByOrderNumber(ctx context.Context, orderNumber string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows := jr.sqlCreateJobByOrderNumber.QueryRowContext(ctx, orderNumber)

	err := rows.Err()
	if err != nil {
		return 0, fmt.Errorf("при вызове CreateJobByOrderNumber(%s) произошла ошибка %w", orderNumber, err)
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
	var err error
	var tx *sql.Tx
	if tx, err = jr.store.Begin(); err != nil {
		logger.LogSugar.Error(err)
		return err
	}
	rows, err := tx.QueryContext(ctx, `update jobs_order set next_run = now() + interval '10 sec', run_cnt = run_cnt + 1, is_work = false where order_number = $1`, orderNumber)
	if err != nil {
		err = errors.Join(err, tx.Rollback())
		return fmt.Errorf("при вызове UpdateJobByOrderNumber(%s) произошла ошибка %w", orderNumber, err)
	}
	err = rows.Err()
	if err != nil {
		err = errors.Join(rows.Err(), tx.Rollback())
		return fmt.Errorf("при вызове UpdateJobByOrderNumber(%s) произошла ошибка %w", orderNumber, err)
	}
	if err = tx.Commit(); err != nil {
		logger.LogSugar.Error(err)
		return nil
	}
	return nil
}

func (jr *JobRepository) DeleteJobByOrderNumber(ctx context.Context, orderNumber string) error {
	ctx, cancel := context.WithTimeout(ctx, config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	var err error
	var tx *sql.Tx
	if tx, err = jr.store.Begin(); err != nil {
		logger.LogSugar.Error(err)
		return err
	}
	rows, err := tx.QueryContext(ctx, `delete from jobs_order where order_number = $1`, orderNumber)
	if err != nil {
		err = errors.Join(err, tx.Rollback())
		return fmt.Errorf("при вызове DeleteJobByOrderNumber(%s) произошла ошибка %w", orderNumber, err)
	}
	err = rows.Err()
	if err != nil {
		err = errors.Join(rows.Err(), tx.Rollback())
		return fmt.Errorf("при вызове DeleteJobByOrderNumber(%s) произошла ошибка %w", orderNumber, err)
	}
	if err = tx.Commit(); err != nil {
		logger.LogSugar.Error(err)
		return nil
	}
	return nil
}
