package utils

import (
	"net/url"
	"testing"
)

func TestParsePositive(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		present bool
		want    int
		err     string
	}{
		{name: "ok", raw: "3", present: true, want: 3},
		{name: "missing", present: false, err: "int1 is required"},
		{name: "zero", raw: "0", present: true, err: "int1 must be greater than 0"},
		{name: "not a number", raw: "no", present: true, err: "int1 must be an integer"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePositive(tt.raw, "int1", tt.present)
			if tt.err != "" {
				if err == nil || err.Error() != tt.err {
					t.Fatalf("err %v", err)
				}
				return
			}
			if err != nil || got != tt.want {
				t.Fatalf("got %d err %v", got, err)
			}
		})
	}
}

func TestParseWord(t *testing.T) {
	q := url.Values{}
	if _, err := ParseWord(q, "str1", 3); err == nil || err.Error() != "str1 is required" {
		t.Fatalf("err %v", err)
	}

	q.Set("str1", "")
	got, err := ParseWord(q, "str1", 3)
	if err != nil || got != "" {
		t.Fatalf("got %q err %v", got, err)
	}

	q.Set("str1", "abcd")
	if _, err := ParseWord(q, "str1", 3); err == nil || err.Error() != "str1 is too long" {
		t.Fatalf("err %v", err)
	}
}
