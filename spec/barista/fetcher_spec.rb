# typed: false
# frozen_string_literal: true

require "barista"

RSpec.describe Barista::Fetcher do
  describe ".fetch" do
    it "returns the response body on success" do
      stub_request(:get, "https://example.com/sla")
        .to_return(status: 200, body: "<html>SLA content</html>")

      result = described_class.fetch("https://example.com/sla")
      expect(result).to eq("<html>SLA content</html>")
    end

    it "raises on non-success responses" do
      stub_request(:get, "https://example.com/missing")
        .to_return(status: 404, body: "Not Found")

      expect { described_class.fetch("https://example.com/missing") }.to raise_error(RuntimeError, /404/)
    end
  end
end
