# typed: strict
# frozen_string_literal: true

require_relative "barista/providers/service"
require_relative "barista/providers/provider"
require_relative "barista/configuration"
require_relative "barista/fetcher"
require_relative "barista/workers/service_guarantee_explorer_worker"

# Top-level module for the Barista project.
module Barista
  # Namespace for Sidekiq job classes.
  module Workers; end

  extend T::Sig

  sig { returns(Configuration) }
  def self.configuration
    @configuration ||= T.let(Configuration.load, T.nilable(Configuration))
  end

  sig { params(config: Configuration).void }
  def self.configure(config)
    @configuration = T.let(config, T.nilable(Configuration))
  end
end
