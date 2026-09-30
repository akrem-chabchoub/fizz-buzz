package memory

import (
	"context"
	"sync"
	"testing"

	"github.com/akrem-chabchoub/fizz-buzz/internal/stats"
)

func TestMemoryMostUsedWins(t *testing.T) {
	s := New()
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

func TestMemoryTiePrefersSmallerInt1(t *testing.T) {
	s := New()
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

func TestMemoryParallelRecords(t *testing.T) {
	s := New()
	req := stats.Request{Int1: 3, Int2: 5, Limit: 10, Str1: "a", Str2: "b"}
	const n = 50
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if err := s.Record(context.Background(), req); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	_, hits, ok, err := s.Top(context.Background())
	if err != nil || !ok {
		t.Fatalf("ok %v err %v", ok, err)
	}
	if hits != n {
		t.Fatalf("hits %d", hits)
	}
}

func TestMemoryCancelledContextIsNotCounted(t *testing.T) {
	s := New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := s.Record(ctx, stats.Request{Int1: 2, Int2: 3, Limit: 1, Str1: "a", Str2: "b"})
	if err == nil {
		t.Fatal("expected error")
	}
	_, _, ok, err := s.Top(context.Background())
	if err != nil || ok {
		t.Fatalf("ok %v err %v", ok, err)
	}
}
