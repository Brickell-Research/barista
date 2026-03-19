# frozen_string_literal: true

require "barista"

RSpec.describe Barista::Providers::Service do
  subject(:service) { described_class.new(provider_name: "aws", name: "s3", url: "https://aws.amazon.com/s3/sla/") }

  it "exposes its attributes" do
    expect(service).to have_attributes(provider_name: "aws", name: "s3", url: "https://aws.amazon.com/s3/sla/")
  end

  describe "#key" do
    it "returns provider/service format" do
      expect(service.key).to eq("aws/s3")
    end
  end

  describe ".parse_key" do
    it "splits a key into provider_name and service_name" do
      expect(described_class.parse_key("aws/s3")).to eq(provider_name: "aws", service_name: "s3")
    end

    it "handles service names containing slashes" do
      result = described_class.parse_key("gcp/cloud/storage")
      expect(result).to eq(provider_name: "gcp", service_name: "cloud/storage")
    end
  end
end
