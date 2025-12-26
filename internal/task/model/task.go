package model

import "time"

type Task struct {
	ID          int64
	Title       string
	Description string
	Status      string
	Assignee    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
