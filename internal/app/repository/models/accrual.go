package models

import "time"

type Accrual struct {
	ID        int       `json:"id"`
	Order     Order     `json:"order"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}
