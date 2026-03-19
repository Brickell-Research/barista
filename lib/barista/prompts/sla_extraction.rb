# frozen_string_literal: true

module Barista
  module Prompts
    module SlaExtraction
      SYSTEM = <<~PROMPT.freeze
        You are an SLA parser. Extract service level availability guarantees from the provided HTML documentation.

        Return ONLY a valid JSON object with this structure:
        {
          "guarantees": [
            {
              "name": "snake_case_name",
              "threshold": 99.9,
              "window_days": 30
            }
          ]
        }

        Rules:
        - Only include guarantees with an explicit percentage uptime or availability threshold
        - name must be a short snake_case identifier (e.g. "monthly_uptime_percentage")
        - threshold is a Float (e.g. 99.9, not "99.9%")
        - window_days is an Integer representing the SLA measurement window in days
        - If no guarantees are found, return {"guarantees": []}
        - Return only the raw JSON object, no preamble, explanation, or code fences
      PROMPT

      def self.user(service:, content:)
        <<~PROMPT
          Extract SLA guarantees from the following documentation.

          Provider: #{service.provider_name}
          Service: #{service.name}
          Source URL: #{service.url}

          Content:
          #{content}
        PROMPT
      end
    end
  end
end
