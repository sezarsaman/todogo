package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"task-manager/internal/observability/metrics"
	"task-manager/internal/task/model"

	"github.com/gin-gonic/gin"
)

// TaskService defines the interface for task service operations
type TaskService interface {
	Get(ctx context.Context, id int64) (*model.Task, error)
	List(ctx context.Context) ([]model.Task, error)
	Create(ctx context.Context, task *model.Task) error
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id int64) error
}

type Handler struct {
	service TaskService
}

func New(svc TaskService) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) Register(r *gin.Engine) {
	r.POST("/tasks", h.create)
	r.GET("/tasks", h.list)
	r.GET("/tasks/:id", h.get)
	r.PUT("/tasks/:id", h.update)
	r.DELETE("/tasks/:id", h.delete)
}

// CreateTaskRequest represents the request payload for creating a task
type CreateTaskRequest struct {
	Title       string `json:"title" example:"Buy groceries"`
	Description string `json:"description" example:"Buy milk, eggs, and bread"`
	Status      string `json:"status" example:"pending" default:"pending"`
	Assignee    string `json:"assignee" example:"john.doe@example.com"`
}

// TaskResponse represents the response payload for task operations
type TaskResponse struct {
	ID          int64  `json:"id" example:"1"`
	Title       string `json:"title" example:"Buy groceries"`
	Description string `json:"description" example:"Buy milk, eggs, and bread"`
	Status      string `json:"status" example:"pending"`
	Assignee    string `json:"assignee" example:"john.doe@example.com"`
	CreatedAt   string `json:"created_at" example:"2024-12-26T10:30:00Z"`
	UpdatedAt   string `json:"updated_at" example:"2024-12-26T10:30:00Z"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"title is required"`
}

// CreateTask godoc
// @Summary      Create a new task
// @Description  Create a new task with title, description, status, and assignee
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body CreateTaskRequest true "Task request"
// @Success      201  {object}  model.Task "Task created successfully"
// @Failure      400  {object}  ErrorResponse "Invalid request"
// @Failure      422  {object}  ErrorResponse "Validation error"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /tasks [post]
func (h *Handler) create(c *gin.Context) {
	var req model.Task

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "title is required"})
		return
	}

	if len(req.Title) < 3 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "title is too short"})
		return
	}

	if err := h.service.Create(c.Request.Context(), &req); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Record task creation metric
	metrics.IncrementTasksCount(req.Status)

	c.JSON(http.StatusCreated, req)
}

// ListTasks godoc
// @Summary      List all tasks
// @Description  Retrieve a list of all tasks
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Success      200  {array}   model.Task "List of tasks"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /tasks [get]
func (h *Handler) list(c *gin.Context) {
	tasks, err := h.service.List(c.Request.Context())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetTask godoc
// @Summary      Get a task by ID
// @Description  Retrieve a specific task by its ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path      int64  true  "Task ID"
// @Success      200  {object}  model.Task "Task details"
// @Failure      400  {object}  ErrorResponse "Invalid task ID"
// @Failure      404  {object}  ErrorResponse "Task not found"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /tasks/{id} [get]
func (h *Handler) get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	task, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask godoc
// @Summary      Update a task
// @Description  Update an existing task by ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id       path  int64  true  "Task ID"
// @Param        request  body  CreateTaskRequest  true  "Updated task data"
// @Success      202  "Task updated successfully"
// @Failure      400  {object}  ErrorResponse "Invalid request"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /tasks/{id} [put]
func (h *Handler) update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var req model.Task
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	req.ID = id
	if err := h.service.Update(c.Request.Context(), &req); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Record task update metric (decrement old status, increment new status)
	// For simplicity, we just record the update with new status
	metrics.IncrementTasksCount(req.Status)

	c.Status(http.StatusAccepted)
}

// DeleteTask godoc
// @Summary      Delete a task
// @Description  Delete a task by its ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id  path  int64  true  "Task ID"
// @Success      204  "Task deleted successfully"
// @Failure      400  {object}  ErrorResponse "Invalid task ID"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /tasks/{id} [delete]
func (h *Handler) delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Record task deletion metric
	metrics.DecrementTasksCount("pending")

	c.Status(http.StatusNoContent)
}
