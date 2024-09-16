package models

import (
	"database/sql"
)

type Order struct {
	ID        int             `json:"id,omitempty"`
	Number    string          `json:"number"`
	Status    string          `json:"status"`
	User      User            `json:"user"`
	Accrual   sql.NullFloat64 `json:"accrual"`
	CreatedAt sql.NullTime    `json:"created_at,omitempty"`
	DeletedAt sql.NullTime    `json:"deleted_at,omitempty"`
}
