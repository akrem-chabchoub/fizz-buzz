package config

import "testing"

func TestLoadRejectsBadShutdownTimeout(t *testing.T) {
	t.Setenv("SHUTDOWN_TIMEOUT", "nope")
	if _, err := Load(); err == nil {
		t.Fatal("expected error")
	}
}
