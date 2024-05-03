package entity

import "time"

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password,omitempty"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
