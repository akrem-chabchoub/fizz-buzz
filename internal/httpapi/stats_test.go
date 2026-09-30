package httpapi

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestStatsEmpty(t *testing.T) {
	status, body := do(t, testHandler(100, 64), "/v1/stats")
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
	h := testHandler(100, 64)
	status, body := do(t, h, "/v1/fizzbuzz?int1=0&int2=5&limit=5&str1=fizz&str2=buzz")
	if status != http.StatusBadRequest {
		t.Fatalf("status %d body %s", status, body)
	}
	status, body = do(t, h, "/v1/stats")
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

func TestStatsMostUsed(t *testing.T) {
	h := testHandler(100, 64)
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
