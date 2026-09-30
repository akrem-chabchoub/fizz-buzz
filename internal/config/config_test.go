package config

import "testing"

func TestLoadRejectsBadShutdownTimeout(t *testing.T) {
	t.Setenv("SHUTDOWN_TIMEOUT", "nope")
	if _, err := Load(); err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadRejectsNonPositiveMaxLimit(t *testing.T) {
	t.Setenv("MAX_LIMIT", "0")
	if _, err := Load(); err == nil {
		t.Fatal("expected error")
	}
}
