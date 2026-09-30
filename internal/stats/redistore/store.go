package redistore

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/akrem-chabchoub/fizz-buzz/internal/stats"
)

const requestsKey = "fizzbuzz:requests"

// topScript returns the member with the highest score.
// Equal scores go to the lexicographically smallest member, same as memory.
var topScript = redis.NewScript(`
local top = redis.call('ZREVRANGE', KEYS[1], 0, 0, 'WITHSCORES')
if top[1] == nil then
  return nil
end
local score = top[2]
local members = redis.call('ZRANGEBYSCORE', KEYS[1], score, score)
table.sort(members)
return {members[1], score}
`)

// Store keeps the tally in Redis.
type Store struct {
	client *redis.Client
}

func New(client *redis.Client) *Store {
	return &Store{client: client}
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
	if err == redis.Nil {
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
	hits, err := scoreHits(pair[1])
	if err != nil {
		return stats.Request{}, 0, false, err
	}
	var req stats.Request
	if err := json.Unmarshal([]byte(member), &req); err != nil {
		return stats.Request{}, 0, false, err
	}
	return req, hits, true, nil
}

func scoreHits(v any) (uint64, error) {
	switch n := v.(type) {
	case int64:
		if n < 0 {
			return 0, fmt.Errorf("stats top: negative score")
		}
		return uint64(n), nil
	case string:
		f, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return 0, fmt.Errorf("stats top: score %q", n)
		}
		if f < 0 {
			return 0, fmt.Errorf("stats top: negative score")
		}
		return uint64(f), nil
	default:
		return 0, fmt.Errorf("stats top: score %T", v)
	}
}
