package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/akrem-chabchoub/fizz-buzz/internal/stats"
)

type errStore struct{}

func (errStore) Record(context.Context, stats.Request) error {
	return errors.New("boom")
}

func (errStore) Top(context.Context) (stats.Request, uint64, bool, error) {
	return stats.Request{}, 0, false, errors.New("boom")
}

func (errStore) Ping(context.Context) error {
	return errors.New("boom")
}

type blockStore struct{}

func (blockStore) Record(ctx context.Context, _ stats.Request) error {
	<-ctx.Done()
	return ctx.Err()
}

func (blockStore) Top(ctx context.Context) (stats.Request, uint64, bool, error) {
	<-ctx.Done()
	return stats.Request{}, 0, false, ctx.Err()
}

func (blockStore) Ping(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func TestListReturnsRowWhenRecordFails(t *testing.T) {
	svc := New(errStore{}, time.Second)
	got := svc.List(context.Background(), stats.Request{
		Int1: 2, Int2: 3, Limit: 1, Str1: "a", Str2: "b",
	})
	if len(got) != 1 || got[0] != "1" {
		t.Fatalf("%#v", got)
	}
}

func TestListReturnsRowWhenRecordTimesOut(t *testing.T) {
	svc := New(blockStore{}, 20*time.Millisecond)
	got := svc.List(context.Background(), stats.Request{
		Int1: 2, Int2: 3, Limit: 1, Str1: "a", Str2: "b",
	})
	if len(got) != 1 || got[0] != "1" {
		t.Fatalf("%#v", got)
	}
}

func TestTopAndReadyStopAtTimeout(t *testing.T) {
	svc := New(blockStore{}, 20*time.Millisecond)
	if _, _, _, err := svc.Top(context.Background()); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("top err %v", err)
	}
	if err := svc.Ready(context.Background()); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("ready err %v", err)
	}
}
