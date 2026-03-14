# frozen_string_literal: true

source "https://rubygems.org"

ruby "4.0.1"

gem "sidekiq", "~> 7.3"
gem "sidekiq-cron", "~> 2.0"
gem "httparty"
gem "sorbet-runtime"

group :development, :test do
  gem "rspec", "~> 3.13"
  gem "rubocop", "~> 1.75", require: false
  gem "rubocop-rspec", "~> 3.5", require: false
  gem "rubocop-sorbet", require: false
  gem "sorbet-static", require: false
  gem "tapioca", require: false
end
