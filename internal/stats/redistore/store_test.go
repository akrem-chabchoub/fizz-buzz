package redistore

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"github.com/akrem-chabchoub/fizz-buzz/internal/stats"
)

func TestRedisMostUsedWins(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	popular := stats.Request{Int1: 3, Int2: 5, Limit: 10, Str1: "fizz", Str2: "buzz"}
	other := stats.Request{Int1: 2, Int2: 5, Limit: 10, Str1: "fizz", Str2: "buzz"}
	if err := s.Record(ctx, popular); err != nil {
		t.Fatal(err)
	}
	if err := s.Record(ctx, popular); err != nil {
		t.Fatal(err)
	}
	if err := s.Record(ctx, other); err != nil {
		t.Fatal(err)
	}
	got, hits, ok, err := s.Top(ctx)
	if err != nil || !ok {
		t.Fatalf("ok %v err %v", ok, err)
	}
	if got != popular || hits != 2 {
		t.Fatalf("got %+v hits %d", got, hits)
	}
}

func TestRedisTiePrefersSmallerInt1(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	higher := stats.Request{Int1: 3, Int2: 5, Limit: 10, Str1: "fizz", Str2: "buzz"}
	lower := stats.Request{Int1: 2, Int2: 5, Limit: 10, Str1: "fizz", Str2: "buzz"}
	if err := s.Record(ctx, higher); err != nil {
		t.Fatal(err)
	}
	if err := s.Record(ctx, lower); err != nil {
		t.Fatal(err)
	}
	got, hits, ok, err := s.Top(ctx)
	if err != nil || !ok {
		t.Fatalf("ok %v err %v", ok, err)
	}
	if got.Int1 != 2 || hits != 1 {
		t.Fatalf("got %+v hits %d", got, hits)
	}
}

func TestRedisEmptyTop(t *testing.T) {
	s := newStore(t)
	_, _, ok, err := s.Top(context.Background())
	if err != nil || ok {
		t.Fatalf("ok %v err %v", ok, err)
	}
}

func TestRedisPing(t *testing.T) {
	s := newStore(t)
	if err := s.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func newStore(t *testing.T) *Store {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	return New(client)
}
