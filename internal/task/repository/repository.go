package repository

import (
	"context"

	"task-manager/internal/task/model"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id int64) (*model.Task, error)
	List(ctx context.Context) ([]model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id int64) error
}
