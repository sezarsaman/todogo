package http

import (
	_ "task-manager/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// TaskHandlerRegister defines the interface for registering task routes
type TaskHandlerRegister interface {
	Register(r *gin.Engine)
}

func NewRouter(taskHandler TaskHandlerRegister) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           86400,
	}))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	// Swagger documentation endpoints
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})

	// Register task routes
	taskHandler.Register(r)

	return r
}
