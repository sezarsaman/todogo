package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"task-manager/internal/config"
	"task-manager/internal/platform/cache"
	"task-manager/internal/platform/database"
	httpPlatform "task-manager/internal/platform/http"
	taskCache "task-manager/internal/task/cache"
	"task-manager/internal/task/handler"
	"task-manager/internal/task/repository"
	"task-manager/internal/task/service"
)

// @title           Task Manager API
// @version         1.0
// @description     A simple task management API with CRUD operations
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @host            localhost:8080
// @BasePath        /
// @schemes         http
func main() {
	cfg := config.Load()

	db, err := database.New(context.Background(), postgresDSN(cfg))
	if err != nil {
		panic(err)
	}

	redisPlatform, err := cache.New(cfg.RedisHost + ":" + cfg.RedisPort)
	if err != nil {
		panic(err)
	}

	taskRedisCache := taskCache.NewRedis(redisPlatform.Client, 30*time.Second)

	taskRepo := repository.NewPostgres(db.Conn)
	taskService := service.New(taskRepo, taskRedisCache)
	taskHandler := handler.New(taskService)

	router := httpPlatform.NewRouter(taskHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	_ = srv.Shutdown(shutdownCtx)
	_ = redisPlatform.Close()
	_ = db.Conn.Close(context.Background())
}

func postgresDSN(cfg config.Config) string {
	return "postgres://" +
		cfg.DBUser + ":" +
		cfg.DBPassword + "@" +
		cfg.DBHost + ":" +
		cfg.DBPort + "/" +
		cfg.DBName
}
