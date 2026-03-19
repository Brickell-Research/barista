# Barista

Automated explorer for discovering and structuring third-party service guarantees for the Caffeine ecosystem.

## Commands

- `make ci` — runs lint and test
- `bundle exec rubocop` — lint
- `bundle exec rspec` — test
- `bundle exec sidekiq -C ./config/sidekiq.yml` — start workers (requires Redis)

## Sidekiq

Use `Sidekiq::Job` (not `Sidekiq::Worker`). Keep `perform` args JSON-serializable (String, Integer, Float, Boolean, nil, Array, Hash).

## Style

- Ruby 3.4, double quotes, 120-char line limit
- `# frozen_string_literal: true` in every file
