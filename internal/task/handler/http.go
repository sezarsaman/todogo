package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

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

	c.JSON(http.StatusCreated, req)
}

func (h *Handler) list(c *gin.Context) {
	tasks, err := h.service.List(c.Request.Context())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

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

	c.Status(http.StatusAccepted)
}

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

	c.Status(http.StatusNoContent)
}
