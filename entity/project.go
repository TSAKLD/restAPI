package entity

import "time"

type Project struct {
	ID        int64     `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	UserID    int64     `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
