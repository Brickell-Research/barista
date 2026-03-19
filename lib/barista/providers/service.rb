# frozen_string_literal: true

module Barista
  module Providers
    # A single service offered by a provider (e.g. S3 under AWS).
    class Service
      attr_reader :provider_name, :name, :url

      def initialize(provider_name:, name:, url:)
        @provider_name = provider_name
        @name = name
        @url = url
      end

      def key = "#{provider_name}/#{name}"

      def self.parse_key(key)
        provider_name, service_name = key.split("/", 2)
        { provider_name:, service_name: }
      end
    end
  end
end
