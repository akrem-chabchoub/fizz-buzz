package memory

import (
	"context"
	"sync"

	"github.com/akrem-chabchoub/fizz-buzz/internal/stats"
)

type entry struct {
	req  stats.Request
	hits uint64
}

// Store keeps the tally in this process.
type Store struct {
	mu sync.Mutex
	m  map[string]entry
}

func New() *Store {
	return &Store{m: map[string]entry{}}
}

// Record adds one hit. A cancelled context is not counted.
func (s *Store) Record(ctx context.Context, req stats.Request) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	key, err := stats.Canonical(req)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e := s.m[key]
	e.req = req
	e.hits++
	s.m[key] = e
	return nil
}

// Top returns the request with the most hits.
// Equal hits go to the smaller canonical JSON.
func (s *Store) Top(ctx context.Context) (stats.Request, uint64, bool, error) {
	if err := ctx.Err(); err != nil {
		return stats.Request{}, 0, false, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var (
		bestKey string
		best    entry
		found   bool
	)
	for key, e := range s.m {
		if !found || e.hits > best.hits || (e.hits == best.hits && key < bestKey) {
			bestKey = key
			best = e
			found = true
		}
	}
	if !found {
		return stats.Request{}, 0, false, nil
	}
	return best.req, best.hits, true, nil
}
