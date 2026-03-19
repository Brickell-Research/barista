# frozen_string_literal: true

require_relative "barista/providers/service"
require_relative "barista/providers/provider"
require_relative "barista/configuration"
require_relative "barista/fetcher"
require_relative "barista/intermediate"
require_relative "barista/prompts/sla_extraction"
require_relative "barista/llm_client"
require_relative "barista/synthesizer"
require_relative "barista/translator"
require_relative "barista/caffeine_writer"
require_relative "barista/workers/service_guarantee_explorer_worker"

# Top-level module for the Barista project.
module Barista
  # Namespace for Sidekiq job classes.
  module Workers; end

  def self.configuration
    @configuration ||= Configuration.load
  end

  def self.configure(config)
    @configuration = config
  end
end
