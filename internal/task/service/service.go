package service

import (
	"context"

	"task-manager/internal/task/model"
	"task-manager/internal/task/repository"
)

type Service struct {
	repo repository.TaskRepository
}

func New(repo repository.TaskRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, task *model.Task) error {
	if task.Status == "" {
		task.Status = "pending"
	}
	return s.repo.Create(ctx, task)
}

func (s *Service) Get(ctx context.Context, id int64) (*model.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]model.Task, error) {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, task *model.Task) error {
	return s.repo.Update(ctx, task)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
