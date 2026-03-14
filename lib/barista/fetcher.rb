# typed: strict
# frozen_string_literal: true

require "httparty"

module Barista
  # Fetches the content of a URL over HTTP(S).
  module Fetcher
    extend T::Sig

    sig { params(url: String).returns(String) }
    def self.fetch(url)
      response = HTTParty.get(url, headers: { "User-Agent" => "Barista/1.0" }, follow_redirects: true, timeout: 10)
      raise "HTTP #{response.code}: #{response.message}" unless response.success?

      T.must(response.body)
    end
  end
end
