package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/akrem-chabchoub/fizz-buzz/internal/app"
	"github.com/akrem-chabchoub/fizz-buzz/internal/stats/memory"
)

func TestReadyNilPing(t *testing.T) {
	status, body := do(t, testHandler(100, 64), "/ready")
	if status != http.StatusOK {
		t.Fatalf("status %d body %s", status, body)
	}
	assertStatus(t, body, "ok")
}

func TestReadyDown(t *testing.T) {
	h := New(app.New(memory.New(), time.Second), 100, 64, func(context.Context) error {
		return errors.New("down")
	}).Handler()
	status, body := do(t, h, "/ready")
	if status != http.StatusServiceUnavailable {
		t.Fatalf("status %d body %s", status, body)
	}
	assertStatus(t, body, "not ready")
}

func assertStatus(t *testing.T, body []byte, want string) {
	t.Helper()
	var got map[string]string
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	if got["status"] != want {
		t.Fatalf("%v", got)
	}
}
