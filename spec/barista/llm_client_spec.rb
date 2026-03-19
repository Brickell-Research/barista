# frozen_string_literal: true

require "barista"

RSpec.describe Barista::LlmClient do
  let(:service) do
    Barista::Providers::Service.new(provider_name: "aws", name: "s3", url: "https://aws.amazon.com/s3/sla/")
  end
  let(:content) { "<html>AWS S3 guarantees 99.9% monthly uptime.</html>" }

  let(:anthropic_response) do
    {
      id: "msg_123",
      type: "message",
      role: "assistant",
      content: [{ type: "text", text: response_text }],
      model: Barista::LlmClient::MODEL,
      stop_reason: "end_turn",
      usage: { input_tokens: 100, output_tokens: 50 }
    }.to_json
  end

  before { ENV["ANTHROPIC_API_KEY"] = "test-key" }
  after  { ENV.delete("ANTHROPIC_API_KEY") }

  describe ".extract_guarantees" do
    context "when the API returns guarantees" do
      let(:response_text) do
        '{"guarantees":[{"name":"monthly_uptime_percentage","threshold":99.9,"window_days":30}]}'
      end

      before do
        stub_request(:post, "https://api.anthropic.com/v1/messages")
          .to_return(status: 200, body: anthropic_response, headers: { "Content-Type" => "application/json" })
      end

      it "returns parsed guarantees" do
        result = described_class.extract_guarantees(service:, content:)
        expect(result["guarantees"].size).to eq(1)
      end

      it "returns the correct guarantee attributes" do
        result = described_class.extract_guarantees(service:, content:)
        expect(result["guarantees"].first).to include(
          "name" => "monthly_uptime_percentage",
          "threshold" => 99.9,
          "window_days" => 30
        )
      end
    end

    context "when the API returns no guarantees" do
      let(:response_text) { '{"guarantees":[]}' }

      before do
        stub_request(:post, "https://api.anthropic.com/v1/messages")
          .to_return(status: 200, body: anthropic_response, headers: { "Content-Type" => "application/json" })
      end

      it "returns an empty guarantees array" do
        result = described_class.extract_guarantees(service:, content:)
        expect(result["guarantees"]).to be_empty
      end
    end

    context "when the API returns a server error" do
      before do
        stub_request(:post, "https://api.anthropic.com/v1/messages")
          .to_return(status: 500, body: '{"error":{"type":"api_error","message":"Internal server error"}}',
            headers: { "Content-Type" => "application/json" })
      end

      it "raises an error" do
        expect { described_class.extract_guarantees(service:, content:) }.to raise_error(StandardError)
      end
    end
  end
end
