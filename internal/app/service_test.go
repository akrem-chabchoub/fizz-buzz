package app

import (
	"context"
	"errors"
	"testing"

	"github.com/akrem-chabchoub/fizz-buzz/internal/stats"
)

type errStore struct{}

func (errStore) Record(context.Context, stats.Request) error {
	return errors.New("boom")
}

func (errStore) Top(context.Context) (stats.Request, uint64, bool, error) {
	return stats.Request{}, 0, false, errors.New("boom")
}

func TestListReturnsRowWhenRecordFails(t *testing.T) {
	svc := New(errStore{})
	got := svc.List(context.Background(), stats.Request{
		Int1: 2, Int2: 3, Limit: 1, Str1: "a", Str2: "b",
	})
	if len(got) != 1 || got[0] != "1" {
		t.Fatalf("%#v", got)
	}
}
