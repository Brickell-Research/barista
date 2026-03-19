# frozen_string_literal: true

require "barista"
require "tmpdir"

RSpec.describe Barista::Configuration do
  describe ".load" do
    let(:yaml) do
      <<~YAML
        providers:
          - name: stripe
            services:
              - name: payments
                url: https://docs.stripe.com/api
          - name: aws
            services:
              - name: s3
                url: https://aws.amazon.com/s3/sla/
              - name: rds
                url: https://aws.amazon.com/rds/sla/
      YAML
    end

    let(:path) do
      p = File.join(Dir.tmpdir, "test_services.yml")
      File.write(p, yaml)
      p
    end

    let(:config) { described_class.load(path) }

    after { FileUtils.rm_f(path) }

    it "loads the correct number of providers" do
      expect(config.providers.size).to eq(2)
    end

    it "parses provider names" do
      expect(config.providers.map(&:name)).to contain_exactly("stripe", "aws")
    end

    it "parses services under each provider" do
      aws = config.providers.find { |p| p.name == "aws" }
      expect(aws.services.size).to eq(2)
    end

    it "sets provider_name on each service" do
      stripe = config.providers.find { |p| p.name == "stripe" }
      expect(stripe.services.first.provider_name).to eq("stripe")
    end
  end

  describe "#output_dir" do
    it "defaults to 'output'" do
      config = described_class.new(providers: [])
      expect(config.output_dir).to eq("output")
    end
  end
end
