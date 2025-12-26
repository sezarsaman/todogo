package cache

import (
	"context"

	"task-manager/internal/task/model"
)

type TaskCache interface {
	GetTasks(ctx context.Context) ([]model.Task, bool, error)
	SetTasks(ctx context.Context, tasks []model.Task) error
	Invalidate(ctx context.Context) error
}
