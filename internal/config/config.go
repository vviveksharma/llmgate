package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr         string
	PostgresDSN      string
	RedisAddr        string
	DefaultProvider  string
	RateLimitRPS     float64
	RateLimitBurst   int
	JWTSecret        string
	ShutdownTimeout  time.Duration
}

func Load() (*Config, error) {
	return &Config{
		HTTPAddr:        getEnv("HTTP_ADDR", ":8080"),
		PostgresDSN:     getEnv("POSTGRES_DSN", ""),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		DefaultProvider: getEnv("DEFAULT_PROVIDER", "mock"),
		RateLimitRPS:    getEnvFloat("RATE_LIMIT_RPS", 100),
		RateLimitBurst:  getEnvInt("RATE_LIMIT_BURST", 200),
		JWTSecret:       getEnv("JWT_SECRET", "dev-secret-change-me"),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
