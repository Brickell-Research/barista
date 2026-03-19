# frozen_string_literal: true

module Barista
  # Calls the LLM to extract SLA guarantees from raw content and returns an Intermediate.
  module Synthesizer
    def self.synthesize(service:, content:)
      raw = LlmClient.extract_guarantees(service:, content:)

      guarantees = Array(raw["guarantees"]).map do |g|
        Guarantee.new(name: g["name"], threshold: g["threshold"].to_f, window_days: g["window_days"].to_i)
      end

      Intermediate.new(
        service_name: service.name,
        provider_name: service.provider_name,
        source_url: service.url,
        guarantees: guarantees
      )
    end
  end
end
