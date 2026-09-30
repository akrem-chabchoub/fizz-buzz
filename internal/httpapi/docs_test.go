package httpapi

import (
	"net/http"
	"strings"
	"testing"
)

func TestOpenAPI(t *testing.T) {
	status, body := do(t, testHandler(100, 64), "/openapi.yaml")
	if status != http.StatusOK {
		t.Fatalf("status %d", status)
	}
	text := string(body)
	if !strings.Contains(text, "/v1/fizzbuzz") || !strings.Contains(text, "/ready") {
		t.Fatalf("%s", text)
	}
}

func TestSwagger(t *testing.T) {
	status, body := do(t, testHandler(100, 64), "/swagger/")
	if status != http.StatusOK {
		t.Fatalf("status %d body %s", status, body)
	}
	if !strings.Contains(strings.ToLower(string(body)), "swagger") {
		t.Fatalf("%s", body)
	}
}
