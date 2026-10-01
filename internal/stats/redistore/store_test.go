package redistore

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	"github.com/akrem-chabchoub/fizz-buzz/internal/stats"
)

func TestRedisMostUsedWins(t *testing.T) {
	s, _ := newStore(t)
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

func TestRedisTieGoesToSmallerJSON(t *testing.T) {
	s, _ := newStore(t)
	ctx := context.Background()
	nine := stats.Request{Int1: 9, Int2: 5, Limit: 10, Str1: "fizz", Str2: "buzz"}
	ten := stats.Request{Int1: 10, Int2: 5, Limit: 10, Str1: "fizz", Str2: "buzz"}
	if err := s.Record(ctx, nine); err != nil {
		t.Fatal(err)
	}
	if err := s.Record(ctx, ten); err != nil {
		t.Fatal(err)
	}
	got, hits, ok, err := s.Top(ctx)
	if err != nil || !ok {
		t.Fatalf("ok %v err %v", ok, err)
	}
	if got != ten || hits != 1 {
		t.Fatalf("got %+v hits %d", got, hits)
	}
}

func TestRedisEmptyTop(t *testing.T) {
	s, _ := newStore(t)
	_, _, ok, err := s.Top(context.Background())
	if err != nil || ok {
		t.Fatalf("ok %v err %v", ok, err)
	}
}

func TestRedisTopRejectsBadMember(t *testing.T) {
	s, mr := newStore(t)
	if _, err := mr.ZAdd(requestsKey, 1, "not json"); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := s.Top(context.Background()); err == nil {
		t.Fatal("expected error")
	}
}

func TestRedisPing(t *testing.T) {
	s, _ := newStore(t)
	if err := s.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestRedisRecordStopsAtDeadline(t *testing.T) {
	s := New(silentAddr(t), "")
	t.Cleanup(func() { _ = s.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	if err := s.Record(ctx, stats.Request{Int1: 1, Int2: 1, Limit: 1}); err == nil {
		t.Fatal("expected error")
	}
	if d := time.Since(start); d > time.Second {
		t.Fatalf("took %s", d)
	}
}

func newStore(t *testing.T) (*Store, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	s := New(mr.Addr(), "")
	t.Cleanup(func() { _ = s.Close() })
	return s, mr
}

// silentAddr accepts connections and never answers, like a hung Redis.
func silentAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	t.Cleanup(func() {
		close(done)
		_ = ln.Close()
	})
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go func() {
				<-done
				_ = c.Close()
			}()
		}
	}()
	return ln.Addr().String()
}
