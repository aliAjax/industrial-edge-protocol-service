package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr       string
	DataDir        string
	WorkerInterval time.Duration
	MaxBufferBytes int64
	CommandTimeout time.Duration
}

func Load() Config {
	return Config{HTTPAddr: env("EDGE_HTTP_ADDR", ":8099"), DataDir: env("EDGE_DATA_DIR", "./data"), WorkerInterval: duration("EDGE_WORKER_INTERVAL", 5*time.Second), MaxBufferBytes: integer("EDGE_MAX_BUFFER_BYTES", 64<<20), CommandTimeout: duration("EDGE_COMMAND_TIMEOUT", 5*time.Second)}
}
func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func duration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
func integer(key string, fallback int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
}
