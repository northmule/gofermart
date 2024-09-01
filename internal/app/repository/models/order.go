package models

import "time"

type Order struct {
	ID        int       `json:"id,omitempty"`
	Number    string    `json:"number"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
