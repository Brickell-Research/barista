# typed: strict
# frozen_string_literal: true

require "fileutils"

module Barista
  # Writes generated .caffeine expectation files to the output directory.
  module CaffeineWriter
    extend T::Sig

    sig { params(service: Providers::Service, content: String).returns(String) }
    def self.write(service:, content:)
      dir = File.join(Barista.configuration.output_dir, "expectations")
      FileUtils.mkdir_p(dir)

      path = File.join(dir, "#{service.name}.caffeine")
      File.write(path, content)
      path
    end
  end
end
