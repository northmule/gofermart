package models

import (
	"database/sql"
	"github.com/shopspring/decimal"
)

type Balance struct {
	ID        int             `json:"id"`
	User      User            `json:"user"`
	Order     Order           `json:"order"`
	Value     decimal.Decimal `json:"value"`
	UpdatedAt sql.NullTime    `json:"updated_at"`
}
