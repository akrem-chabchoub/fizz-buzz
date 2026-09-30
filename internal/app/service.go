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
	store         stats.Store
	recordTimeout time.Duration
}

func New(store stats.Store, recordTimeout time.Duration) *Service {
	return &Service{store: store, recordTimeout: recordTimeout}
}

// List builds the fizzbuzz list and counts the request.
// The counter waits at most recordTimeout. The list is still returned when that fails.
func (s *Service) List(ctx context.Context, req stats.Request) []string {
	recCtx, cancel := context.WithTimeout(ctx, s.recordTimeout)
	defer cancel()
	if err := s.store.Record(recCtx, req); err != nil {
		slog.Error("stats record failed", "err", err)
	}
	return fizzbuzz.Generate(req.Int1, req.Int2, req.Limit, req.Str1, req.Str2)
}

// Top returns the most used request.
func (s *Service) Top(ctx context.Context) (stats.Request, uint64, bool, error) {
	return s.store.Top(ctx)
}
