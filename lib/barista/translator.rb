# frozen_string_literal: true

module Barista
  # Translates an Intermediate into a valid Caffeine v6 expectations string.
  module Translator
    def self.translate(intermediate)
      return empty_file(intermediate) if intermediate.guarantees.empty?

      lines = ["Expectations"]

      intermediate.guarantees.each do |guarantee|
        lines << ""
        lines << "\"#{guarantee.name}\":"
        lines << "  Guarantees #{format_threshold(guarantee.threshold)}% over #{guarantee.window_days}d window"
      end

      lines << ""
      lines.join("\n")
    end

    def self.format_threshold(threshold)
      int = threshold.to_i
      threshold == int ? int : threshold
    end
    private_class_method :format_threshold

    def self.empty_file(intermediate)
      "# No guarantees found for #{intermediate.provider_name}/#{intermediate.service_name}\n" \
        "# Source: #{intermediate.source_url}\n"
    end
    private_class_method :empty_file
  end
end
