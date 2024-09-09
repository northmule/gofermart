package models

import "database/sql"

type Job struct {
	ID          int          `json:"id,omitempty"`
	OrderNumber string       `json:"order_number,omitempty"`
	NextRun     sql.NullTime `json:"next_run,omitempty"`
	RunCnt      int          `json:"run_cnt,omitempty"`
	CreatedAt   sql.NullTime `json:"created_at,omitempty"`
	UpdatedAt   sql.NullTime `json:"updated_at,omitempty"`
}
