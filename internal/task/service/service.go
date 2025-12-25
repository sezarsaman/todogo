package service

import (
	"context"

	"task-manager/internal/task/cache"
	"task-manager/internal/task/model"
	"task-manager/internal/task/repository"
)

type Service struct {
	repo  repository.TaskRepository
	cache cache.TaskCache
}

func New(repo repository.TaskRepository, cache cache.TaskCache) *Service {
	return &Service{repo: repo, cache: cache}
}

func (s *Service) Get(ctx context.Context, id int64) (*model.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]model.Task, error) {
	if s.cache != nil {
		if tasks, hit, err := s.cache.GetTasks(ctx); err == nil && hit {
			return tasks, nil
		}
	}

	tasks, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		_ = s.cache.SetTasks(ctx, tasks)
	}

	return tasks, nil
}

func (s *Service) Create(ctx context.Context, task *model.Task) error {
	if task.Status == "" {
		task.Status = "pending"
	}
	if err := s.repo.Create(ctx, task); err != nil {
		return err
	}
	if s.cache != nil {
		_ = s.cache.Invalidate(ctx)
	}
	return nil
}

func (s *Service) Update(ctx context.Context, task *model.Task) error {
	if err := s.repo.Update(ctx, task); err != nil {
		return err
	}
	if s.cache != nil {
		_ = s.cache.Invalidate(ctx)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	if s.cache != nil {
		_ = s.cache.Invalidate(ctx)
	}
	return nil
}
