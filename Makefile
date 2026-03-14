.PHONY: ci lint typecheck test

ci: lint typecheck test

lint:
	bundle exec rubocop

typecheck:
	bundle exec srb tc

test:
	bundle exec rspec
