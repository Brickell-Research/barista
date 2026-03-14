# typed: strict
# frozen_string_literal: true

require "yaml"

module Barista
  # Loads the service registry from config/services.yml and exposes settings.
  class Configuration < T::Struct
    const :providers, T::Array[Providers::Provider]
    const :output_dir, String, default: "output"

    extend T::Sig

    sig { params(path: String).returns(Configuration) }
    def self.load(path = File.expand_path("../../config/services.yml", __dir__))
      yaml = YAML.safe_load_file(path, permitted_classes: [], symbolize_names: true)
      providers = Array(yaml[:providers]).map { |p| parse_provider(p) }

      new(providers: providers)
    end

    sig { params(hash: T::Hash[Symbol, T.untyped]).returns(Providers::Provider) }
    def self.parse_provider(hash)
      provider_name = hash[:name]
      services = Array(hash[:services]).map do |s|
        Providers::Service.new(provider_name: provider_name, name: s[:name], url: s[:url])
      end

      Providers::Provider.new(name: hash[:name], services: services)
    end
    private_class_method :parse_provider
  end
end
