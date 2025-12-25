package http

import (
	"task-manager/internal/task/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(taskHandler *handler.Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	taskHandler.Register(r)

	return r
}
