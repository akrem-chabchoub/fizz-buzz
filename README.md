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
- If two questions have the same number of hits, the one with the smaller JSON wins.
- If Redis blips during a call, you still get the list. `/ready` is what fails.
- The set `fizzbuzz:requests` does not expire. `docker compose down` deletes the container, and the tally goes with it.
