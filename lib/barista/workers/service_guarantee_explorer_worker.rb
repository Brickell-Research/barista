# frozen_string_literal: true

require "sidekiq"

module Barista
  module Workers
    # Discovers third-party services and explores their published guarantees.
    #
    # Fan-out pattern:
    #   - Called without args: discovers services and enqueues individual jobs
    #   - Called with a "provider/service" key: explores that specific service
    class ServiceGuaranteeExplorerWorker
      include Sidekiq::Job

      sidekiq_options queue: "exploration"

      def perform(service_key = nil)
        if service_key
          explore_service(service_key)
        else
          discover_and_enqueue_services
        end
      end

      private

      def explore_service(service_key)
        parsed = Providers::Service.parse_key(service_key)
        service = find_service(parsed[:provider_name], parsed[:service_name])

        unless service
          logger.warn("Unknown service key: #{service_key}")
          return
        end

        logger.info("[#{service_key}] Fetching SLA from #{service.url}")
        content = Fetcher.fetch(service.url)

        logger.info("[#{service_key}] Synthesizing guarantees (#{content.length} bytes)")
        intermediate = Synthesizer.synthesize(service:, content:)

        logger.info("[#{service_key}] Found #{intermediate.guarantees.size} guarantee(s)")
        caffeine = Translator.translate(intermediate)

        path = CaffeineWriter.write(service:, content: caffeine)
        logger.info("[#{service_key}] Wrote #{path}")
      end

      def find_service(provider_name, service_name)
        provider = Barista.configuration.providers.find { |p| p.name == provider_name }
        provider&.services&.find { |s| s.name == service_name }
      end

      def discover_and_enqueue_services
        keys = discovered_services
        logger.info("Enqueuing #{keys.size} service exploration job(s)")
        keys.each { |key| self.class.perform_async(key) }
      end

      def discovered_services
        Barista.configuration.providers.flat_map do |provider|
          provider.services.map(&:key)
        end
      end
    end
  end
end
