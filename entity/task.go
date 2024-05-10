package entity

import "time"

type Task struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserID      int64     `json:"user_id"`
	ProjectID   int64     `json:"project_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskToCreate struct {
	Name        string `json:"name"`
	ProjectID   int64  `json:"project_id"`
	Description string `json:"description"`
}
