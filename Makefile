.PHONY: test test-race lint run

test:
	go test ./...

test-race:
	go test -race ./...

lint:
	golangci-lint run ./...

run:
	go run ./cmd/api