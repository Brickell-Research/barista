# frozen_string_literal: true

module Barista
  # Translates an Intermediate into a valid Caffeine Unmeasured Expectations string.
  module Translator
    def self.translate(intermediate)
      return empty_file(intermediate) if intermediate.guarantees.empty?

      lines = ["Unmeasured Expectations"]

      intermediate.guarantees.each do |guarantee|
        lines << ""
        lines << "  * \"#{guarantee.name}\":"
        lines << "    Provides {"
        lines << "      threshold: #{format_threshold(guarantee.threshold)},"
        lines << "      window_in_days: #{guarantee.window_days}"
        lines << "    }"
      end

      lines << ""
      lines.join("\n")
    end

    def self.format_threshold(threshold)
      int = threshold.to_i
      "#{threshold == int ? int : threshold}%"
    end
    private_class_method :format_threshold

    def self.empty_file(intermediate)
      "# No guarantees found for #{intermediate.provider_name}/#{intermediate.service_name}\n" \
        "# Source: #{intermediate.source_url}\n"
    end
    private_class_method :empty_file
  end
end
