package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatsEmpty(t *testing.T) {
	h, _ := testHandler(t, 100, 64)
	status, body := do(t, h, "/v1/stats")
	if status != http.StatusOK {
		t.Fatalf("status %d body %s", status, body)
	}
	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got["hits"] != float64(0) {
		t.Fatalf("%v", got)
	}
}

func TestStatsIgnoresInvalidCall(t *testing.T) {
	h, _ := testHandler(t, 100, 64)
	status, body := do(t, h, "/v1/fizzbuzz?int1=0&int2=5&limit=5&str1=fizz&str2=buzz")
	if status != http.StatusBadRequest {
		t.Fatalf("status %d body %s", status, body)
	}
	assertNoHits(t, h)
}

func TestStatsIgnoresHead(t *testing.T) {
	h, _ := testHandler(t, 100, 64)
	req := httptest.NewRequest(http.MethodHead, "/v1/fizzbuzz?int1=3&int2=5&limit=5&str1=fizz&str2=buzz", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || rec.Body.Len() != 0 {
		t.Fatalf("status %d body %s", rec.Code, rec.Body)
	}
	assertNoHits(t, h)
}

func TestStatsMostUsed(t *testing.T) {
	h, _ := testHandler(t, 100, 64)
	popular := "/v1/fizzbuzz?int1=3&int2=5&limit=5&str1=fizz&str2=buzz"
	other := "/v1/fizzbuzz?int1=2&int2=5&limit=5&str1=fizz&str2=buzz"
	for _, path := range []string{popular, popular, other} {
		status, body := do(t, h, path)
		if status != http.StatusOK {
			t.Fatalf("status %d body %s", status, body)
		}
	}
	status, body := do(t, h, "/v1/stats")
	if status != http.StatusOK {
		t.Fatalf("status %d body %s", status, body)
	}
	var got statsBody
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	if got.Int1 != 3 || got.Int2 != 5 || got.Limit != 5 || got.Str1 != "fizz" || got.Str2 != "buzz" || got.Hits != 2 {
		t.Fatalf("%+v", got)
	}
}

func TestStatsRedisDown(t *testing.T) {
	h, mr := testHandler(t, 100, 64)
	mr.Close()
	status, body := do(t, h, "/v1/stats")
	if status != http.StatusInternalServerError {
		t.Fatalf("status %d body %s", status, body)
	}
}

func assertNoHits(t *testing.T, h http.Handler) {
	t.Helper()
	status, body := do(t, h, "/v1/stats")
	if status != http.StatusOK {
		t.Fatalf("status %d body %s", status, body)
	}
	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	if got["hits"] != float64(0) {
		t.Fatalf("%v", got)
	}
}
