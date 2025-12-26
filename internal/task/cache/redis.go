package cache

import (
	"context"
	"encoding/json"
	"time"

	"task-manager/internal/task/model"

	"github.com/redis/go-redis/v9"
)

const tasksKey = "tasks:all"

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedis(client *redis.Client, ttl time.Duration) *RedisCache {
	return &RedisCache{client: client, ttl: ttl}
}

func (r *RedisCache) GetTasks(ctx context.Context) ([]model.Task, bool, error) {
	val, err := r.client.Get(ctx, tasksKey).Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var tasks []model.Task
	if err := json.Unmarshal([]byte(val), &tasks); err != nil {
		return nil, false, err
	}

	return tasks, true, nil
}

func (r *RedisCache) SetTasks(ctx context.Context, tasks []model.Task) error {
	b, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, tasksKey, b, r.ttl).Err()
}

func (r *RedisCache) Invalidate(ctx context.Context) error {
	return r.client.Del(ctx, tasksKey).Err()
}
