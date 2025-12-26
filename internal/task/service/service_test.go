package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"task-manager/internal/task/model"
)

func TestServiceCreate(t *testing.T) {
	ctx := context.Background()
	repo := &MockRepository{}
	cache := &MockCache{}
	svc := New(repo, cache)

	t.Run("create_with_default_status", func(t *testing.T) {
		repo.CreateFunc = func(ctx context.Context, task *model.Task) error {
			task.ID = 1
			task.CreatedAt = time.Now()
			task.UpdatedAt = time.Now()
			return nil
		}
		cache.GetTasksFunc = func(ctx context.Context) ([]model.Task, bool, error) {
			return []model.Task{{ID: 1, Title: "existing"}}, true, nil
		}
		cache.SetTasksFunc = func(ctx context.Context, tasks []model.Task) error {
			if len(tasks) != 2 {
				t.Errorf("expected 2 tasks, got %d", len(tasks))
			}
			return nil
		}

		task := &model.Task{Title: "New Task"}
		err := svc.Create(ctx, task)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if task.Status != "pending" {
			t.Errorf("expected status pending, got %s", task.Status)
		}
	})

	t.Run("create_with_repo_error", func(t *testing.T) {
		repo.CreateFunc = func(ctx context.Context, task *model.Task) error {
			return errors.New("db error")
		}
		task := &model.Task{Title: "New Task"}
		err := svc.Create(ctx, task)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}

func TestServiceUpdate(t *testing.T) {
	ctx := context.Background()
	repo := &MockRepository{}
	cache := &MockCache{}
	svc := New(repo, cache)

	t.Run("update_with_cache_hit", func(t *testing.T) {
		repo.UpdateFunc = func(ctx context.Context, task *model.Task) error {
			return nil
		}
		cache.GetTasksFunc = func(ctx context.Context) ([]model.Task, bool, error) {
			return []model.Task{
				{ID: 1, Title: "old title", Status: "pending"},
				{ID: 2, Title: "other"},
			}, true, nil
		}
		cache.SetTasksFunc = func(ctx context.Context, tasks []model.Task) error {
			if tasks[0].Title != "updated title" {
				t.Errorf("expected updated title, got %s", tasks[0].Title)
			}
			return nil
		}

		task := &model.Task{ID: 1, Title: "updated title", Status: "in_progress"}
		err := svc.Update(ctx, task)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("update_with_cache_miss", func(t *testing.T) {
		repo.UpdateFunc = func(ctx context.Context, task *model.Task) error {
			return nil
		}
		repo.ListFunc = func(ctx context.Context) ([]model.Task, error) {
			return []model.Task{{ID: 1, Title: "task1"}}, nil
		}
		cache.GetTasksFunc = func(ctx context.Context) ([]model.Task, bool, error) {
			return nil, false, nil
		}
		cache.SetTasksFunc = func(ctx context.Context, tasks []model.Task) error {
			return nil
		}

		task := &model.Task{ID: 1, Title: "updated"}
		err := svc.Update(ctx, task)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestServiceDelete(t *testing.T) {
	ctx := context.Background()
	repo := &MockRepository{}
	cache := &MockCache{}
	svc := New(repo, cache)

	t.Run("delete_with_cache_hit", func(t *testing.T) {
		repo.DeleteFunc = func(ctx context.Context, id int64) error {
			return nil
		}
		cache.GetTasksFunc = func(ctx context.Context) ([]model.Task, bool, error) {
			return []model.Task{
				{ID: 1, Title: "task1"},
				{ID: 2, Title: "task2"},
			}, true, nil
		}
		cache.SetTasksFunc = func(ctx context.Context, tasks []model.Task) error {
			if len(tasks) != 1 {
				t.Errorf("expected 1 task after delete, got %d", len(tasks))
			}
			return nil
		}

		err := svc.Delete(ctx, 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("delete_with_repo_error", func(t *testing.T) {
		repo.DeleteFunc = func(ctx context.Context, id int64) error {
			return errors.New("db error")
		}
		err := svc.Delete(ctx, 1)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("delete_with_cache_error", func(t *testing.T) {
		repo.DeleteFunc = func(ctx context.Context, id int64) error {
			return nil
		}
		cache.GetTasksFunc = func(ctx context.Context) ([]model.Task, bool, error) {
			return nil, false, errors.New("cache error")
		}
		repo.ListFunc = func(ctx context.Context) ([]model.Task, error) {
			return []model.Task{{ID: 1, Title: "task"}}, nil
		}
		cache.SetTasksFunc = func(ctx context.Context, tasks []model.Task) error {
			return nil
		}

		err := svc.Delete(ctx, 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestServiceGet(t *testing.T) {
	ctx := context.Background()
	repo := &MockRepository{}
	cache := &MockCache{}
	svc := New(repo, cache)

	t.Run("get_from_db_and_sync_cache", func(t *testing.T) {
		repo.GetByIDFunc = func(ctx context.Context, id int64) (*model.Task, error) {
			return &model.Task{ID: 1, Title: "task1", Status: "pending"}, nil
		}
		cache.GetTasksFunc = func(ctx context.Context) ([]model.Task, bool, error) {
			return []model.Task{}, true, nil
		}
		cache.SetTasksFunc = func(ctx context.Context, tasks []model.Task) error {
			return nil
		}

		task, err := svc.Get(ctx, 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if task.Title != "task1" {
			t.Errorf("expected task1, got %s", task.Title)
		}
	})

	t.Run("get_with_repo_error", func(t *testing.T) {
		repo.GetByIDFunc = func(ctx context.Context, id int64) (*model.Task, error) {
			return nil, errors.New("not found")
		}
		_, err := svc.Get(ctx, 999)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("get_with_cache_refresh_on_miss", func(t *testing.T) {
		repo.GetByIDFunc = func(ctx context.Context, id int64) (*model.Task, error) {
			return &model.Task{ID: 1, Title: "task1"}, nil
		}
		cache.GetTasksFunc = func(ctx context.Context) ([]model.Task, bool, error) {
			return nil, false, nil
		}
		repo.ListFunc = func(ctx context.Context) ([]model.Task, error) {
			return []model.Task{{ID: 1, Title: "task1"}}, nil
		}
		cache.SetTasksFunc = func(ctx context.Context, tasks []model.Task) error {
			return nil
		}

		task, err := svc.Get(ctx, 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if task.ID != 1 {
			t.Errorf("expected task id 1, got %d", task.ID)
		}
	})
}

func TestServiceList(t *testing.T) {
	ctx := context.Background()
	repo := &MockRepository{}
	cache := &MockCache{}
	svc := New(repo, cache)

	t.Run("list_with_cache_hit", func(t *testing.T) {
		cachedTasks := []model.Task{
			{ID: 1, Title: "task1"},
			{ID: 2, Title: "task2"},
		}
		cache.GetTasksFunc = func(ctx context.Context) ([]model.Task, bool, error) {
			return cachedTasks, true, nil
		}

		tasks, err := svc.List(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 2 {
			t.Errorf("expected 2 tasks, got %d", len(tasks))
		}
	})

	t.Run("list_with_cache_miss", func(t *testing.T) {
		repoTasks := []model.Task{
			{ID: 1, Title: "task1"},
			{ID: 2, Title: "task2"},
		}
		repo.ListFunc = func(ctx context.Context) ([]model.Task, error) {
			return repoTasks, nil
		}
		cache.GetTasksFunc = func(ctx context.Context) ([]model.Task, bool, error) {
			return nil, false, nil
		}
		cache.SetTasksFunc = func(ctx context.Context, tasks []model.Task) error {
			return nil
		}

		tasks, err := svc.List(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 2 {
			t.Errorf("expected 2 tasks, got %d", len(tasks))
		}
	})

	t.Run("list_with_repo_error", func(t *testing.T) {
		repo.ListFunc = func(ctx context.Context) ([]model.Task, error) {
			return nil, errors.New("db error")
		}
		cache.GetTasksFunc = func(ctx context.Context) ([]model.Task, bool, error) {
			return nil, false, nil
		}

		_, err := svc.List(ctx)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}
