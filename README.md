# fizzbuzz

Counts from 1 to a limit. Multiples of `int1` become `str1`, multiples of `int2` become `str2`, and multiples of both become the two words stuck together.

## run

```bash
docker compose up --build
```

Then open http://localhost:8080/swagger/

```bash
make test
```

To hack on the code, start Redis and the server that reloads when you save a file:

```bash
docker compose up -d redis
make dev
```

Settings and their defaults are in `.env.exp`.

## calls

```bash
curl -s localhost:8080/health
curl -s localhost:8080/ready
curl -s "localhost:8080/v1/fizzbuzz?int1=3&int2=5&limit=15&str1=fizz&str2=buzz"
curl -s localhost:8080/v1/stats
```

`/health` means this process is up. `/ready` means Redis answered. The fizzbuzz call returns one string per number, from 1 to `limit`. `/v1/stats` is the question people asked most. Nothing counted yet looks like `{"hits":0}`.

## notes

- When `int1` and `int2` are the same number, every multiple becomes the two words stuck together.
- If two questions have the same number of hits, the one with the smaller JSON wins. So `int1=10` beats `int1=9`.
- Only a GET that returns the list is counted. A bad query or a `HEAD` call is not.
- If Redis is slow or down, you still get the list, at most `STATS_TIMEOUT` later. `/ready` answers 503 and `/v1/stats` answers 500 within the same time.
- The set `fizzbuzz:requests` does not expire. `docker compose down` deletes the container, and the tally goes with it.

## design choices

- `internal/fizzbuzz` is only the counting logic, no HTTP and no Redis. `internal/httpapi` reads the query and writes JSON. `internal/app` builds the list and counts it. `internal/stats/redistore` is the only code that talks to Redis.
- The tally lives in Redis so it survives a restart and stays the same when more than one copy of the server runs. A map in memory would give each copy its own answer.
- It is one sorted set. A hit is a single `ZINCRBY`, so two calls at the same time never lose a count. Finding the most used request is one small script that stays fast even with many ties.
- The member is the JSON of the five parameters, so `int1=3`, `int1=03` and `int1=+3` count as the same request.
- A call is counted before it is answered. Counts are exact, and each call pays one Redis round trip, capped by `STATS_TIMEOUT`. Counting in the background would be faster but could drop hits when the server stops.
- The list is built in memory, so `MAX_LIMIT` (10000) and `MAX_STRING_LEN` (32 characters) are what keep one call small. The worst case is about 4 MB, because JSON writes `<` as `\u003c`. Raising them costs memory on every big call. Streaming the response would be the next step if bigger lists are needed.
- Every new set of parameters adds a member to `fizzbuzz:requests`, and nothing removes them. 500k different requests take about 70 MB in Redis. That is fine here. In a real deployment, put rate limiting in front of the server and watch Redis memory, or trim the rarely used members, which turns the top into an estimate.
