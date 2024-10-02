package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Withdrawn struct {
	ID        int             `json:"id"`
	UserID    int             `json:"user_id"`
	OrderID   int             `json:"order_id"`
	Value     decimal.Decimal `json:"value"`
	CreatedAt time.Time       `json:"created_at"`
	Order     *Order          `json:"order"`
}
