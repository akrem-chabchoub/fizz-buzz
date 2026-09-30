package httpapi

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestFizzbuzz(t *testing.T) {
	h := testHandler(100, 64)
	short := testHandler(100, 3)
	tests := []struct {
		name   string
		path   string
		status int
		want   []string
		err    string
		h      http.Handler
	}{
		{
			name:   "ok",
			path:   "/v1/fizzbuzz?int1=3&int2=5&limit=5&str1=fizz&str2=buzz",
			status: http.StatusOK,
			want:   []string{"1", "2", "fizz", "4", "buzz"},
			h:      h,
		},
		{
			name:   "missing",
			path:   "/v1/fizzbuzz?int2=5&limit=5&str1=fizz&str2=buzz",
			status: http.StatusBadRequest,
			err:    "int1 is required",
			h:      h,
		},
		{
			name:   "zero",
			path:   "/v1/fizzbuzz?int1=0&int2=5&limit=5&str1=fizz&str2=buzz",
			status: http.StatusBadRequest,
			err:    "int1 must be greater than 0",
			h:      h,
		},
		{
			name:   "not a number",
			path:   "/v1/fizzbuzz?int1=no&int2=5&limit=5&str1=fizz&str2=buzz",
			status: http.StatusBadRequest,
			err:    "int1 must be an integer",
			h:      h,
		},
		{
			name:   "limit too big",
			path:   "/v1/fizzbuzz?int1=3&int2=5&limit=101&str1=a&str2=b",
			status: http.StatusBadRequest,
			err:    "limit is too big",
			h:      h,
		},
		{
			name:   "word too long",
			path:   "/v1/fizzbuzz?int1=3&int2=5&limit=5&str1=abcd&str2=b",
			status: http.StatusBadRequest,
			err:    "str1 is too long",
			h:      short,
		},
		{
			name:   "empty word",
			path:   "/v1/fizzbuzz?int1=2&int2=5&limit=2&str1=&str2=b",
			status: http.StatusOK,
			want:   []string{"1", ""},
			h:      h,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body := do(t, tt.h, tt.path)
			if status != tt.status {
				t.Fatalf("status %d, body %s", status, body)
			}
			if tt.err != "" {
				var got map[string]string
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatal(err)
				}
				if got["error"] != tt.err {
					t.Fatalf("error %q", got["error"])
				}
				return
			}
			var got []string
			if err := json.Unmarshal(body, &got); err != nil {
				t.Fatal(err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %#v", got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("got %#v", got)
				}
			}
		})
	}
}
