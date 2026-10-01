package redistore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/akrem-chabchoub/fizz-buzz/internal/stats"
)

const requestsKey = "fizzbuzz:requests"

// topScript returns the member with the highest score.
// Redis keeps equal scores sorted byte by byte, so LIMIT 0 1 is the smallest member.
var topScript = redis.NewScript(`
local top = redis.call('ZREVRANGE', KEYS[1], 0, 0, 'WITHSCORES')
if top[1] == nil then
  return nil
end
local members = redis.call('ZRANGEBYSCORE', KEYS[1], top[2], top[2], 'LIMIT', 0, 1)
return {members[1], tonumber(top[2])}
`)

// Store keeps the tally in Redis.
type Store struct {
	client *redis.Client
}

// New connects to Redis on first use.
// Every call stops at the deadline of its context.
func New(addr, password string) *Store {
	return &Store{client: redis.NewClient(&redis.Options{
		Addr:                  addr,
		Password:              password,
		ContextTimeoutEnabled: true,
	})}
}

// Close releases the connections.
func (s *Store) Close() error {
	return s.client.Close()
}

// Ping reports whether Redis answers.
func (s *Store) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

// Record adds one hit. A cancelled context is not counted.
func (s *Store) Record(ctx context.Context, req stats.Request) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	member, err := stats.Canonical(req)
	if err != nil {
		return err
	}
	return s.client.ZIncrBy(ctx, requestsKey, 1, member).Err()
}

// Top returns the request with the most hits.
// Equal hits go to the smaller canonical JSON.
func (s *Store) Top(ctx context.Context) (stats.Request, uint64, bool, error) {
	if err := ctx.Err(); err != nil {
		return stats.Request{}, 0, false, err
	}
	raw, err := topScript.Run(ctx, s.client, []string{requestsKey}).Result()
	if errors.Is(err, redis.Nil) {
		return stats.Request{}, 0, false, nil
	}
	if err != nil {
		return stats.Request{}, 0, false, err
	}
	pair, ok := raw.([]any)
	if !ok || len(pair) != 2 {
		return stats.Request{}, 0, false, fmt.Errorf("stats top: unexpected reply")
	}
	member, ok := pair[0].(string)
	if !ok {
		return stats.Request{}, 0, false, fmt.Errorf("stats top: member %T", pair[0])
	}
	hits, ok := pair[1].(int64)
	if !ok || hits < 0 {
		return stats.Request{}, 0, false, fmt.Errorf("stats top: score %v", pair[1])
	}
	var req stats.Request
	if err := json.Unmarshal([]byte(member), &req); err != nil {
		return stats.Request{}, 0, false, err
	}
	return req, uint64(hits), true, nil
}
