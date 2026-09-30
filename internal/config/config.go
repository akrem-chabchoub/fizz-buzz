package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr        string
	LogLevel        string
	ShutdownTimeout time.Duration
	// MaxLimit is the largest list the server will build.
	MaxLimit int
	// MaxStringLen is the longest word, counted in runes.
	MaxStringLen int
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:        env("HTTP_ADDR", ":8080"),
		LogLevel:        env("LOG_LEVEL", "info"),
		ShutdownTimeout: 10 * time.Second,
	}
	if raw, ok := os.LookupEnv("SHUTDOWN_TIMEOUT"); ok && raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil || d <= 0 {
			return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT must be a positive duration")
		}
		cfg.ShutdownTimeout = d
	}
	var err error
	if cfg.MaxLimit, err = envInt("MAX_LIMIT", 100_000); err != nil {
		return Config{}, err
	}
	if cfg.MaxStringLen, err = envInt("MAX_STRING_LEN", 64); err != nil {
		return Config{}, err
	}
	if cfg.MaxLimit < 1 {
		return Config{}, fmt.Errorf("MAX_LIMIT must be greater than 0")
	}
	if cfg.MaxStringLen < 1 {
		return Config{}, fmt.Errorf("MAX_STRING_LEN must be greater than 0")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

// envInt reads a whole number. A missing variable keeps fallback.
func envInt(key string, fallback int) (int, error) {
	raw, ok := os.LookupEnv(key)
	if !ok || raw == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", key)
	}
	return n, nil
}
