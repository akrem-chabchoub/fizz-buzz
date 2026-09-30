.PHONY: test test-race run

test:
	go test ./...

test-race:
	go test -race ./...

run:
	go run ./cmd/api