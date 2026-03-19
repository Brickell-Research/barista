# barista

An automated explorer for discovering and structuring third-party service guarantees for the [Caffeine](https://caffeine-lang.run) ecosystem.

***

Every service you depend on publishes guarantees. Barista makes them machine-readable.

## Development

**Requires Go 1.23+**.

```sh
go test ./...
```

## Running Workers

Start Redis, then boot the worker:

```sh
redis-server &
go run ./cmd/barista
```

Connects to `REDIS_URL` (defaults to `redis://localhost:6379/0`). Runs an hourly cron to discover and explore all configured services.

## Local Execution

```sh
ANTHROPIC_API_KEY=_____________ go run ./cmd/explore aws/s3
```

## License

GPL-3.0
