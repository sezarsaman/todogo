package model

import "time"

// Task represents a task item in the system
type Task struct {
	ID          int64     `json:"id" example:"1"`
	Title       string    `json:"title" binding:"required" example:"Buy groceries"`
	Description string    `json:"description" example:"Buy milk, eggs, and bread"`
	Status      string    `json:"status" example:"pending" default:"pending"`
	Assignee    string    `json:"assignee" example:"john.doe@example.com"`
	CreatedAt   time.Time `json:"created_at" example:"2024-12-26T10:30:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-12-26T10:30:00Z"`
}
