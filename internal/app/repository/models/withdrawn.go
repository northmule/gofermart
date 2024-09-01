package models

import "time"

type Withdrawn struct {
	ID          int       `json:"id,omitempty"`
	OrderNumber string    `json:"order_number"`
	Value       float64   `json:"value"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}
