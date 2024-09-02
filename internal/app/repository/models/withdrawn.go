package models

import "time"

type Withdrawn struct {
	ID        int       `json:"id,omitempty"`
	OrderID   string    `json:"order_id"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
