# Barista

Automated explorer for discovering and structuring third-party service guarantees for the Caffeine ecosystem.

## Commands

- `make ci` — runs lint and test
- `go vet ./...` — lint
- `go test ./...` — test
- `go run ./cmd/barista` — start worker (requires Redis)
- `go run ./cmd/explore [provider/service]` — run exploration locally without Redis

## Worker Queue

Uses [asynq](https://github.com/hibiken/asynq) (Redis-backed). Two task types:
- `barista:discover` — fan-out job; enqueues one `barista:explore` per configured service
- `barista:explore` — fetches docs, calls LLM, writes `.caffeine` output file

Task payloads must be JSON-serializable structs.

## Style

- Go 1.23
- Module path: `barista`
- Standard library `log/slog` for structured logging
