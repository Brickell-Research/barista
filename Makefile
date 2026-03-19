.PHONY: ci lint test

ci: lint test

lint:
	bundle exec rubocop

test:
	bundle exec rspec
