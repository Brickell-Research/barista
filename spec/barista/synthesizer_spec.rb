# frozen_string_literal: true

require "barista"

RSpec.describe Barista::Synthesizer do
  describe ".synthesize" do
    subject(:result) do
      described_class.synthesize(service: service, content: "some SLA content")
    end

    let(:service) do
      Barista::Providers::Service.new(provider_name: "aws", name: "s3", url: "https://aws.amazon.com/s3/sla/")
    end

    it "includes the service name in the expectation" do
      expect(result).to include('expectation "s3"')
    end

    it "includes the source URL" do
      expect(result).to include("https://aws.amazon.com/s3/sla/")
    end
  end
end
