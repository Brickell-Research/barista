FROM ruby:3.4.1-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
  build-essential \
  && rm -rf /var/lib/apt/lists/*

COPY Gemfile Gemfile.lock ./
RUN bundle config set --local without "development test" \
  && bundle install --jobs 4 --retry 3

COPY . .

RUN mkdir -p output/expectations

CMD ["bundle", "exec", "sidekiq", "-C", "./config/sidekiq.yml"]
