# typed: strict
# frozen_string_literal: true

module Barista
  # Synthesizes fetched content into a .caffeine expectation string.
  # Currently stubbed — real LLM implementation comes later.
  module Synthesizer
    extend T::Sig

    sig { params(service: Providers::Service, content: String).returns(String) }
    def self.synthesize(service:, content:)
      <<~CAFFEINE
        expectation "#{service.name}" {
          # Auto-generated from #{service.url}
          # Content length: #{content.length} bytes
        }
      CAFFEINE
    end
  end
end
