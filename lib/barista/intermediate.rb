# frozen_string_literal: true

module Barista
  # A single SLA guarantee extracted from a provider's documentation.
  Guarantee = Data.define(:name, :threshold, :window_days)

  # The structured result of LLM synthesis for a single service —
  # a list of guarantees ready to be translated into Caffeine.
  Intermediate = Data.define(:service_name, :provider_name, :source_url, :guarantees)
end
