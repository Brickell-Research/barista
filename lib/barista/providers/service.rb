# typed: strict
# frozen_string_literal: true

require "sorbet-runtime"

module Barista
  module Providers
    # A single service offered by a provider (e.g. S3 under AWS).
    class Service < T::Struct
      const :provider_name, String
      const :name, String
      const :url, String

      extend T::Sig

      sig { returns(String) }
      def key
        "#{provider_name}/#{name}"
      end

      sig { params(key: String).returns({ provider_name: String, service_name: String }) }
      def self.parse_key(key)
        provider_name, service_name = key.split("/", 2)
        { provider_name: T.must(provider_name), service_name: T.must(service_name) }
      end
    end
  end
end
