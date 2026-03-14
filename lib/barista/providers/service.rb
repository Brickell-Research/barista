# typed: strict
# frozen_string_literal: true

module Barista
  module Providers
    # A single service offered by a provider (e.g. S3 under AWS).
    class Service < T::Struct
      const :name, String
      const :url, String
    end
  end
end
