package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akrem-chabchoub/fizz-buzz/internal/app"
	"github.com/akrem-chabchoub/fizz-buzz/internal/stats/memory"
)

func testHandler(maxLimit, maxStringLen int) http.Handler {
	return New(app.New(memory.New(), time.Second), maxLimit, maxStringLen, nil).Handler()
}

func TestHealth(t *testing.T) {
	status, body := do(t, testHandler(100, 64), "/health")
	if status != http.StatusOK {
		t.Fatalf("status %d", status)
	}
	var got map[string]string
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	if got["status"] != "ok" {
		t.Fatalf("%v", got)
	}
}

func do(t *testing.T, h http.Handler, path string) (int, []byte) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, rec.Body.Bytes()
}
