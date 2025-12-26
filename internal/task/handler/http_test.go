package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"task-manager/internal/task/model"

	"github.com/gin-gonic/gin"
)

// MockService is a mock of the TaskService interface
type MockService struct {
	GetFunc    func(ctx context.Context, id int64) (*model.Task, error)
	ListFunc   func(ctx context.Context) ([]model.Task, error)
	CreateFunc func(ctx context.Context, task *model.Task) error
	UpdateFunc func(ctx context.Context, task *model.Task) error
	DeleteFunc func(ctx context.Context, id int64) error
}

func (m *MockService) Get(ctx context.Context, id int64) (*model.Task, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockService) List(ctx context.Context) ([]model.Task, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}

func (m *MockService) Create(ctx context.Context, task *model.Task) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, task)
	}
	return nil
}

func (m *MockService) Update(ctx context.Context, task *model.Task) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, task)
	}
	return nil
}

func (m *MockService) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func TestHandlerCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("create_valid_task", func(t *testing.T) {
		mockSvc := &MockService{
			CreateFunc: func(ctx context.Context, task *model.Task) error {
				task.ID = 1
				task.CreatedAt = time.Now()
				task.UpdatedAt = time.Now()
				return nil
			},
		}
		handler := New(mockSvc)
		engine := gin.New()
		handler.Register(engine)

		body := map[string]interface{}{
			"Title":       "New Task",
			"Description": "Task description",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		engine.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}
	})

	t.Run("create_missing_title", func(t *testing.T) {
		mockSvc := &MockService{}
		handler := New(mockSvc)
		engine := gin.New()
		handler.Register(engine)

		body := map[string]interface{}{
			"Description": "Task description",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		engine.ServeHTTP(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", w.Code)
		}
	})

	t.Run("create_title_too_short", func(t *testing.T) {
		mockSvc := &MockService{}
		handler := New(mockSvc)
		engine := gin.New()
		handler.Register(engine)

		body := map[string]interface{}{
			"Title": "ab",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		engine.ServeHTTP(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", w.Code)
		}
	})
}

func TestHandlerList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("list_tasks", func(t *testing.T) {
		mockSvc := &MockService{
			ListFunc: func(ctx context.Context) ([]model.Task, error) {
				return []model.Task{
					{ID: 1, Title: "Task 1", Status: "pending"},
					{ID: 2, Title: "Task 2", Status: "done"},
				}, nil
			},
		}
		handler := New(mockSvc)
		engine := gin.New()
		handler.Register(engine)

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()

		engine.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var tasks []model.Task
		json.Unmarshal(w.Body.Bytes(), &tasks)
		if len(tasks) != 2 {
			t.Errorf("expected 2 tasks, got %d", len(tasks))
		}
	})
}

func TestHandlerGet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("get_task", func(t *testing.T) {
		mockSvc := &MockService{
			GetFunc: func(ctx context.Context, id int64) (*model.Task, error) {
				return &model.Task{ID: 1, Title: "Task 1", Status: "pending"}, nil
			},
		}
		handler := New(mockSvc)
		engine := gin.New()
		handler.Register(engine)

		req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
		w := httptest.NewRecorder()

		engine.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var task model.Task
		json.Unmarshal(w.Body.Bytes(), &task)
		if task.ID != 1 {
			t.Errorf("expected task id 1, got %d", task.ID)
		}
	})

	t.Run("get_task_not_found", func(t *testing.T) {
		mockSvc := &MockService{
			GetFunc: func(ctx context.Context, id int64) (*model.Task, error) {
				// Return error to simulate not found
				return nil, errors.New("not found")
			},
		}
		handler := New(mockSvc)
		engine := gin.New()
		handler.Register(engine)

		req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
		w := httptest.NewRecorder()

		engine.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})

	t.Run("get_task_invalid_id", func(t *testing.T) {
		mockSvc := &MockService{}
		handler := New(mockSvc)
		engine := gin.New()
		handler.Register(engine)

		req := httptest.NewRequest(http.MethodGet, "/tasks/abc", nil)
		w := httptest.NewRecorder()

		engine.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestHandlerUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("update_task", func(t *testing.T) {
		mockSvc := &MockService{
			UpdateFunc: func(ctx context.Context, task *model.Task) error {
				return nil
			},
		}
		handler := New(mockSvc)
		engine := gin.New()
		handler.Register(engine)

		body := map[string]interface{}{
			"Title":  "Updated Task",
			"Status": "done",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		engine.ServeHTTP(w, req)

		if w.Code != http.StatusAccepted {
			t.Errorf("expected status 202, got %d", w.Code)
		}
	})
}

func TestHandlerDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("delete_task", func(t *testing.T) {
		mockSvc := &MockService{
			DeleteFunc: func(ctx context.Context, id int64) error {
				return nil
			},
		}
		handler := New(mockSvc)
		engine := gin.New()
		handler.Register(engine)

		req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
		w := httptest.NewRecorder()

		engine.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", w.Code)
		}
	})
}
