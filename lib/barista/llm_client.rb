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
      JSON.parse(strip_code_fences(response.content.first.text))
    end

    def self.strip_code_fences(text)
      text.strip.sub(/\A```(?:json)?\n?/, "").sub(/\n?```\z/, "")
    end
    private_class_method :strip_code_fences

    def self.build_client
      Anthropic::Client.new
    end
    private_class_method :build_client
  end
end
