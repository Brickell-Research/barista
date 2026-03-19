# frozen_string_literal: true

module Barista
  module Providers
    # A third-party provider (e.g. AWS, Stripe) containing one or more services.
    class Provider
      attr_reader :name, :services

      def initialize(name:, services:)
        @name = name
        @services = services
      end
    end
  end
end
