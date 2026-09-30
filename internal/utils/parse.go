package utils

import (
	"fmt"
	"net/url"
	"strconv"
	"unicode/utf8"
)

// ParsePositive reads a whole number greater than zero.
// present is false when the parameter was not sent.
func ParsePositive(raw, name string, present bool) (int, error) {
	if !present {
		return 0, fmt.Errorf("%s is required", name)
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", name)
	}
	if n <= 0 {
		return 0, fmt.Errorf("%s must be greater than 0", name)
	}
	return n, nil
}

// ParseWord reads a query value.
// An empty value is allowed. A longer value, counted in runes, is not.
func ParseWord(q url.Values, name string, maxLen int) (string, error) {
	if !q.Has(name) {
		return "", fmt.Errorf("%s is required", name)
	}
	v := q.Get(name)
	if utf8.RuneCountInString(v) > maxLen {
		return "", fmt.Errorf("%s is too long", name)
	}
	return v, nil
}
