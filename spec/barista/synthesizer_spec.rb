# frozen_string_literal: true

require "barista"

RSpec.describe Barista::Synthesizer do
  let(:service) do
    Barista::Providers::Service.new(provider_name: "aws", name: "s3", url: "https://aws.amazon.com/s3/sla/")
  end
  let(:content) { "<html>S3 SLA content</html>" }

  let(:llm_response) do
    {
      "guarantees" => [
        { "name" => "monthly_uptime_percentage", "threshold" => 99.9, "window_days" => 30 },
        { "name" => "reduced_redundancy_storage", "threshold" => 99.0, "window_days" => 30 }
      ]
    }
  end

  before do
    allow(Barista::LlmClient).to receive(:extract_guarantees).with(service:, content:).and_return(llm_response)
  end

  describe ".synthesize" do
    subject(:intermediate) { described_class.synthesize(service:, content:) }

    it "returns an Intermediate" do
      expect(intermediate).to be_a(Barista::Intermediate)
    end

    it "sets service metadata" do
      expect(intermediate).to have_attributes(
        service_name: "s3",
        provider_name: "aws",
        source_url: "https://aws.amazon.com/s3/sla/"
      )
    end

    it "builds the correct number of guarantees" do
      expect(intermediate.guarantees.size).to eq(2)
    end

    it "maps guarantee attributes correctly" do
      expect(intermediate.guarantees.first).to have_attributes(
        name: "monthly_uptime_percentage",
        threshold: 99.9,
        window_days: 30
      )
    end

    context "when the LLM returns no guarantees" do
      let(:llm_response) { { "guarantees" => [] } }

      it "returns an Intermediate with an empty guarantees list" do
        expect(intermediate.guarantees).to be_empty
      end
    end
  end
end
