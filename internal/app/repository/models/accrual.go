package models

import "time"

type Accrual struct {
<<<<<<< HEAD
	ID        int       `json:"id,omitempty"`
	Order     Order     `json:"order"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"created_at,omitempty"`
=======
	ID          int       `json:"id,omitempty"`
	OrderNumber string    `json:"order_number"`
	Value       float64   `json:"value"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
>>>>>>> 94746e2 (базовая структура)
}
