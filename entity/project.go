package entity

import "time"

type Project struct {
	ID        int64
	Name      string
	UserID    int64
	CreatedAt time.Time
}
