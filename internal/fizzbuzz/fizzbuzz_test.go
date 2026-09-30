package fizzbuzz

import (
	"reflect"
	"testing"
)

func TestGenerateClassic(t *testing.T) {
	got := Generate(3, 5, 15, "fizz", "buzz")
	want := []string{
		"1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz",
		"buzz", "11", "fizz", "13", "14", "fizzbuzz",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v", got)
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name       string
		int1, int2 int
		limit      int
		str1, str2 string
		want       []string
	}{
		{name: "limit one", int1: 3, int2: 5, limit: 1, str1: "fizz", str2: "buzz", want: []string{"1"}},
		{name: "same divisors", int1: 3, int2: 3, limit: 6, str1: "a", str2: "b", want: []string{"1", "2", "ab", "4", "5", "ab"}},
		{name: "both on last", int1: 2, int2: 3, limit: 6, str1: "a", str2: "b", want: []string{"1", "a", "b", "a", "5", "ab"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Generate(tt.int1, tt.int2, tt.limit, tt.str1, tt.str2)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v", got)
			}
		})
	}
}
