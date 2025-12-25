package main

import (
	"context"
	"time"

	"task-manager/internal/config"
	"task-manager/internal/platform/cache"
	"task-manager/internal/platform/database"
	httpPlatform "task-manager/internal/platform/http"
	taskCache "task-manager/internal/task/cache"
	"task-manager/internal/task/handler"
	"task-manager/internal/task/repository"
	"task-manager/internal/task/service"
	"task-manager/pkg/shutdown"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.New(ctx, postgresDSN(cfg))
	if err != nil {
		panic(err)
	}
	defer db.Conn.Close(ctx)

	redisPlatform, err := cache.New(cfg.RedisHost + ":" + cfg.RedisPort)
	if err != nil {
		panic(err)
	}

	taskRedisCache := taskCache.NewRedis(redisPlatform.Client, 30*time.Second)

	taskRepo := repository.NewPostgres(db.Conn)
	taskService := service.New(taskRepo, taskRedisCache)
	taskHandler := handler.New(taskService)

	router := httpPlatform.NewRouter(taskHandler)

	go shutdown.Wait(ctx, cancel)

	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		panic(err)
	}
}

func postgresDSN(cfg config.Config) string {
	return "postgres://" +
		cfg.DBUser + ":" +
		cfg.DBPassword + "@" +
		cfg.DBHost + ":" +
		cfg.DBPort + "/" +
		cfg.DBName
}
