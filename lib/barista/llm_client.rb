# frozen_string_literal: true

require "anthropic"
require "json"

module Barista
  # Calls the Anthropic API to extract structured SLA guarantees from raw HTML.
  module LlmClient
    MODEL = "claude-opus-4-6"

    def self.extract_guarantees(service:, content:)
      response = build_client.messages.create(
        model: MODEL,
        max_tokens: 1024,
        system: Prompts::SlaExtraction::SYSTEM,
        messages: [{ role: "user", content: Prompts::SlaExtraction.user(service:, content:) }]
      )
      JSON.parse(response.content.first.text)
    end

    def self.build_client
      Anthropic::Client.new
    end
    private_class_method :build_client
  end
end
