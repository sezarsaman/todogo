package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("load_config_with_env_vars", func(t *testing.T) {
		// Set environment variables
		os.Setenv("APP_ENV", "test")
		os.Setenv("HTTP_PORT", "8080")
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpass")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("REDIS_HOST", "redis")
		os.Setenv("REDIS_PORT", "6379")

		cfg := Load()

		if cfg.Env != "test" {
			t.Errorf("expected APP_ENV 'test', got '%s'", cfg.Env)
		}
		if cfg.HTTPPort != "8080" {
			t.Errorf("expected HTTP_PORT '8080', got '%s'", cfg.HTTPPort)
		}
		if cfg.DBHost != "localhost" {
			t.Errorf("expected DB_HOST 'localhost', got '%s'", cfg.DBHost)
		}
		if cfg.DBUser != "testuser" {
			t.Errorf("expected DB_USER 'testuser', got '%s'", cfg.DBUser)
		}
		if cfg.RedisHost != "redis" {
			t.Errorf("expected REDIS_HOST 'redis', got '%s'", cfg.RedisHost)
		}
	})

	t.Run("load_config_missing_env_var_panics", func(t *testing.T) {
		// Unset required variables
		os.Unsetenv("APP_ENV")
		os.Setenv("HTTP_PORT", "8080")
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpass")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("REDIS_HOST", "redis")
		os.Setenv("REDIS_PORT", "6379")

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic when APP_ENV is missing")
			}
		}()

		Load()
	})
}
