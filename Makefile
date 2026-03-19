.PHONY: ci build test lint

ci: lint test

build:
	go build ./...

test:
	go test ./...

lint:
	go vet ./...

run/worker:
	go run ./cmd/barista

run/explore:
	go run ./cmd/explore $(SERVICE)
