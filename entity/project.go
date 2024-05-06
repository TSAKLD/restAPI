package entity

import "time"

type Project struct {
	ID         int64
	Name       string
	OwnerID    int64
	OwnerEmail string
	OwnerName  string
	CreatedAt  time.Time
}
