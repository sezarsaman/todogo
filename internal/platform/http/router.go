package http

import (
	"github.com/gin-gonic/gin"
)

// TaskHandlerRegister defines the interface for registering task routes
type TaskHandlerRegister interface {
	Register(r *gin.Engine)
}

func NewRouter(taskHandler TaskHandlerRegister) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	taskHandler.Register(r)

	return r
}
