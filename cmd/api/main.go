package main

import (
	"context"

	"task-manager/internal/config"
	"task-manager/internal/platform/cache"
	"task-manager/internal/platform/database"
	httpPlatform "task-manager/internal/platform/http"
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

	_, err = cache.New(cfg.RedisHost + ":" + cfg.RedisPort)
	if err != nil {
		panic(err)
	}

	router := httpPlatform.NewRouter()

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
