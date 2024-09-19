package models

import (
	"database/sql"
	"time"
)

type Order struct {
	ID        int             `json:"id"`
	Number    string          `json:"number"`
	Status    string          `json:"status"`
	User      *User           `json:"user"`
	Accrual   sql.NullFloat64 `json:"accrual"`
	CreatedAt time.Time       `json:"created_at"`
	DeletedAt sql.NullTime    `json:"deleted_at"`
}
