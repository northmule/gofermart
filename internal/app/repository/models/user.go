package models

import "time"

type User struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UUID      string    `json:"uuid"`
}
