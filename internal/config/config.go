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
	// RedisAddr is where the tally lives. Default is localhost:6379.
	RedisAddr string
	// RedisPassword is empty when Redis has no password.
	RedisPassword string
	// StatsTimeout is how long one call to the tally may take, /ready included.
	StatsTimeout time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:      env("HTTP_ADDR", ":8080"),
		LogLevel:      env("LOG_LEVEL", "info"),
		RedisAddr:     env("REDIS_ADDR", "localhost:6379"),
		RedisPassword: env("REDIS_PASSWORD", ""),
	}
	var err error
	if cfg.ShutdownTimeout, err = envDuration("SHUTDOWN_TIMEOUT", 10*time.Second); err != nil {
		return Config{}, err
	}
	if cfg.StatsTimeout, err = envDuration("STATS_TIMEOUT", 200*time.Millisecond); err != nil {
		return Config{}, err
	}
	if cfg.MaxLimit, err = envInt("MAX_LIMIT", 10_000); err != nil {
		return Config{}, err
	}
	if cfg.MaxStringLen, err = envInt("MAX_STRING_LEN", 32); err != nil {
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

func envDuration(key string, fallback time.Duration) (time.Duration, error) {
	raw, ok := os.LookupEnv(key)
	if !ok || raw == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(raw)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration", key)
	}
	return d, nil
}
