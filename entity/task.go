package entity

import "time"

type Task struct {
	ID        int64     `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	UserID    int64     `json:"user_id,omitempty"`
	ProjectID int64     `json:"project_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskToCreate struct {
	Name      string `json:"name,omitempty"`
	ProjectID int64  `json:"project_id,omitempty"`
	UserID    int64
}
