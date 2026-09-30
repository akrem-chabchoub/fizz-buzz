package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewRejectsBadLevel(t *testing.T) {
	_, err := New(&bytes.Buffer{}, "loud")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewFiltersDebug(t *testing.T) {
	var buf bytes.Buffer
	log, err := New(&buf, "info")
	if err != nil {
		t.Fatal(err)
	}
	log.Debug("hidden")
	log.Info("shown")
	got := buf.String()
	if strings.Contains(got, "hidden") || !strings.Contains(got, "shown") {
		t.Fatalf("%s", got)
	}
}
