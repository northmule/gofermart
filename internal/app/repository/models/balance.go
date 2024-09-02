package models

import "time"

type Balance struct {
	ID        int       `json:"id,omitempty"`
	OrderID   string    `json:"order_id"`
	Value     float64   `json:"value"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
