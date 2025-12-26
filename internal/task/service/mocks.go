package service

import (
	"context"

	"task-manager/internal/task/model"
)

// MockRepository is a mock of TaskRepository for testing
type MockRepository struct {
	GetByIDFunc func(ctx context.Context, id int64) (*model.Task, error)
	ListFunc    func(ctx context.Context) ([]model.Task, error)
	CreateFunc  func(ctx context.Context, task *model.Task) error
	UpdateFunc  func(ctx context.Context, task *model.Task) error
	DeleteFunc  func(ctx context.Context, id int64) error
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockRepository) List(ctx context.Context) ([]model.Task, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}

func (m *MockRepository) Create(ctx context.Context, task *model.Task) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, task)
	}
	return nil
}

func (m *MockRepository) Update(ctx context.Context, task *model.Task) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, task)
	}
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

// MockCache is a mock of TaskCache for testing
type MockCache struct {
	GetTasksFunc   func(ctx context.Context) ([]model.Task, bool, error)
	SetTasksFunc   func(ctx context.Context, tasks []model.Task) error
	InvalidateFunc func(ctx context.Context) error
}

func (m *MockCache) GetTasks(ctx context.Context) ([]model.Task, bool, error) {
	if m.GetTasksFunc != nil {
		return m.GetTasksFunc(ctx)
	}
	return nil, false, nil
}

func (m *MockCache) SetTasks(ctx context.Context, tasks []model.Task) error {
	if m.SetTasksFunc != nil {
		return m.SetTasksFunc(ctx, tasks)
	}
	return nil
}

func (m *MockCache) Invalidate(ctx context.Context) error {
	if m.InvalidateFunc != nil {
		return m.InvalidateFunc(ctx)
	}
	return nil
}
