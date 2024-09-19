package models

import (
	"database/sql"
)

type Balance struct {
	ID        int          `json:"id"`
	User      User         `json:"user"`
	Order     Order        `json:"order"`
	Value     float64      `json:"value"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
