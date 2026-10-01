package httpapi

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestReadyOK(t *testing.T) {
	h, _ := testHandler(t, 100, 64)
	status, body := do(t, h, "/ready")
	if status != http.StatusOK {
		t.Fatalf("status %d body %s", status, body)
	}
	assertStatus(t, body, "ok")
}

func TestReadyDown(t *testing.T) {
	h, mr := testHandler(t, 100, 64)
	mr.Close()
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
