package config

import "os"

type Config struct {
	Env        string
	HTTPPort   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisHost  string
	RedisPort  string
}

func Load() Config {
	return Config{
		Env:        must("APP_ENV"),
		HTTPPort:   must("HTTP_PORT"),
		DBHost:     must("DB_HOST"),
		DBPort:     must("DB_PORT"),
		DBUser:     must("DB_USER"),
		DBPassword: must("DB_PASSWORD"),
		DBName:     must("DB_NAME"),
		RedisHost:  must("REDIS_HOST"),
		RedisPort:  must("REDIS_PORT"),
	}
}

func must(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("missing env var: " + key)
	}
	return val
}
