# typed: strict
# frozen_string_literal: true

module Barista
  module Providers
    # A third-party provider (e.g. AWS, Stripe) containing one or more services.
    class Provider < T::Struct
      const :name, String
      const :services, T::Array[Service]
    end
  end
end
