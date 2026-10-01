package stats

import (
	"context"
	"encoding/json"
)

// Request is one successful fizzbuzz call.
// Its JSON is the key stored in Redis and breaks ties. Changing a field,
// its order or its tag starts a new tally.
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
	Ping(ctx context.Context) error
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
