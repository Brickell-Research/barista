# typed: strict
# frozen_string_literal: true

require "sorbet-runtime"
require "sidekiq"

module Barista
  module Workers
    # Discovers third-party services and explores their published guarantees.
    #
    # Fan-out pattern:
    #   - Called without args: discovers services and enqueues individual jobs
    #   - Called with a service name: explores that specific service
    class ServiceGuaranteeExplorerWorker
      extend T::Sig
      include Sidekiq::Job

      sidekiq_options queue: "exploration"

      sig { params(service_name: T.nilable(String)).void }
      def perform(service_name = nil)
        if service_name
          explore_service(service_name)
        else
          discover_and_enqueue_services
        end
      end

      private

      sig { params(service_name: String).void }
      def explore_service(service_name)
        # TODO: Implement service-specific guarantee exploration
      end

      sig { void }
      def discover_and_enqueue_services
        discovered_services.each do |service|
          self.class.perform_async(service)
        end
      end

      sig { returns(T::Array[String]) }
      def discovered_services
        # TODO: Implement service discovery
        []
      end
    end
  end
end
