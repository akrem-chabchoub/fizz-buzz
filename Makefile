.PHONY: test test-race lint run dev

test:
	go test ./...

test-race:
	go test -race ./...

lint:
	golangci-lint run ./...

run:
	go run ./cmd/api

dev:
	go run github.com/air-verse/air@v1.62.0