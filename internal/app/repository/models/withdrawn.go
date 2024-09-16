package models

<<<<<<< HEAD
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
=======
import "time"

type Withdrawn struct {
	ID          int       `json:"id,omitempty"`
	OrderNumber string    `json:"order_number"`
	Value       float64   `json:"value"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
>>>>>>> 94746e2 (базовая структура)
}
