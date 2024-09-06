package models

import (
	"database/sql"
)

type Withdrawn struct {
	ID        int          `json:"id,omitempty"`
	UserID    int          `json:"user_id"`
	OrderID   int          `json:"order_id"`
	Value     float64      `json:"value"`
	CreatedAt sql.NullTime `json:"created_at"`
	Order     Order
}
