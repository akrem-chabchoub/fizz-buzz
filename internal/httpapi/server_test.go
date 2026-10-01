package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	"github.com/akrem-chabchoub/fizz-buzz/internal/app"
	"github.com/akrem-chabchoub/fizz-buzz/internal/stats/redistore"
)

func testHandler(t *testing.T, maxLimit, maxStringLen int) (http.Handler, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	store := redistore.New(mr.Addr(), "")
	t.Cleanup(func() { _ = store.Close() })
	return New(app.New(store, time.Second), maxLimit, maxStringLen).Handler(), mr
}

func TestHealth(t *testing.T) {
	h, _ := testHandler(t, 100, 64)
	status, body := do(t, h, "/health")
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

func TestPanicAnswers500(t *testing.T) {
	h := withLog(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	}))
	status, body := do(t, h, "/")
	if status != http.StatusInternalServerError {
		t.Fatalf("status %d body %s", status, body)
	}
	var got map[string]string
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	if got["error"] != "internal error" {
		t.Fatalf("%v", got)
	}
}

func TestPanicAfterWriteKeepsBody(t *testing.T) {
	h := withLog(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, []string{"1"})
		panic("boom")
	}))
	status, body := do(t, h, "/")
	if status != http.StatusOK || strings.TrimSpace(string(body)) != `["1"]` {
		t.Fatalf("status %d body %s", status, body)
	}
}

func do(t *testing.T, h http.Handler, path string) (int, []byte) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, rec.Body.Bytes()
}
