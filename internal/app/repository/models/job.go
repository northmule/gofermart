package models

import (
	"database/sql"
	"time"
)

type Job struct {
	ID          int          `json:"id"`
	OrderNumber string       `json:"order_number"`
	NextRun     sql.NullTime `json:"next_run"`
	RunCnt      int          `json:"run_cnt"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
}
