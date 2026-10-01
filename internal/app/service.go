package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/akrem-chabchoub/fizz-buzz/internal/fizzbuzz"
	"github.com/akrem-chabchoub/fizz-buzz/internal/stats"
)

// Service runs the fizzbuzz flow and the tally.
type Service struct {
	store   stats.Store
	timeout time.Duration
}

// New bounds every store call by timeout.
func New(store stats.Store, timeout time.Duration) *Service {
	return &Service{store: store, timeout: timeout}
}

// List builds the fizzbuzz list and counts the request.
// A failed count is logged and the list is still returned.
func (s *Service) List(ctx context.Context, req stats.Request) []string {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	if err := s.store.Record(ctx, req); err != nil {
		slog.Error("stats record failed", "err", err)
	}
	return fizzbuzz.Generate(req.Int1, req.Int2, req.Limit, req.Str1, req.Str2)
}

// Top returns the most used request.
func (s *Service) Top(ctx context.Context) (stats.Request, uint64, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.store.Top(ctx)
}

// Ready reports whether the store answers.
func (s *Service) Ready(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.store.Ping(ctx)
}
