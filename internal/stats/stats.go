package stats

import (
	"context"
	"encoding/json"
)

// Request is one successful fizzbuzz call.
// Field order is the tie key. Do not reorder the fields.
type Request struct {
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Limit int    `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}

// Store counts requests and returns the most used one.
type Store interface {
	Record(ctx context.Context, req Request) error
	Top(ctx context.Context) (Request, uint64, bool, error)
}

// Canonical is the JSON of the request.
// When two requests have the same number of hits, the smaller string wins.
func Canonical(req Request) (string, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
