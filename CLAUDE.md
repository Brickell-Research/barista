# Barista

Automated explorer for discovering and structuring third-party service guarantees for the Caffeine ecosystem.

## Commands

- `make ci` — runs lint, typecheck, and test
- `bundle exec rubocop` — lint
- `bundle exec srb tc` — type check
- `bundle exec rspec` — test
- `bundle exec sidekiq -C ./config/sidekiq.yml` — start workers (requires Redis)

## Sorbet

Every file in `lib/` must be `# typed: strict`. Specs use `# typed: false`.

### Sigs

All methods in `# typed: strict` files need sigs. Use `extend T::Sig` in every class/module.

```ruby
sig { params(name: String).returns(String) }
def greet(name) = "Hello #{name}"

sig { returns(String) }
attr_reader :name

sig { void }
def reset! = nil
```

### Type assertions (prefer refactoring over these)

- `T.let(expr, Type)` — annotate variable type
- `T.must(expr)` — assert non-nil
- `T.cast(expr, Type)` — when you know more than Sorbet
- `T.unsafe(expr)` — last resort escape hatch; avoid

### RBI files

Use Tapioca for gem type stubs. Hand-written shims go in `sorbet/rbi/shims/`.

- `bin/tapioca gems` — regenerate gem RBIs
- `bin/tapioca dsl` — regenerate DSL RBIs (e.g., Sidekiq workers)
- `bin/tapioca annotations` — pull community type annotations

### Sidekiq workers

Use `Sidekiq::Job` (not `Sidekiq::Worker`). Keep `perform` args JSON-serializable (String, Integer, Float, Boolean, nil, Array, Hash).

## Style

- Ruby 4.0, double quotes, 120-char line limit
- `# frozen_string_literal: true` in every file
