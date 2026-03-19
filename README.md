# barista

An automated explorer for discovering and structuring third-party service guarantees for the [Caffeine](https://caffeine-lang.run) ecosystem.

***

Every service you depend on publishes guarantees. Barista makes them machine-readable.

## Development

**Requires Ruby [4.0](https://www.ruby-lang.org/en/news/2025/12/25/ruby-4-0-0-released/)**.

```sh
bundle exec rspec
```

## Running Workers

Start Redis, then boot Sidekiq:

```sh
redis-server &
bundle exec sidekiq -C ./config/sidekiq.yml
```

Sidekiq connects to `REDIS_URL` (defaults to `redis://localhost:6379/0`). Cron schedules are loaded automatically on startup from `config/schedule.yml`.

## Local Execution

```
ANTHROPIC_API_KEY=_____________ bundle exec ruby bin/explore aws/s3
```

## License

GPL-3.0
