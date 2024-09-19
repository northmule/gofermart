package models

import (
	"time"
)

type Withdrawn struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	OrderID   int       `json:"order_id"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	Order     *Order    `json:"order"`
}
