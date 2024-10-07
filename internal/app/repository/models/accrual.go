package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Accrual struct {
	ID        int             `json:"id"`
	Order     Order           `json:"order"`
	Value     decimal.Decimal `json:"value"`
	CreatedAt time.Time       `json:"created_at"`
}
