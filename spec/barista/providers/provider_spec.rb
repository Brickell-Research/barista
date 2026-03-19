# frozen_string_literal: true

require "barista"

RSpec.describe Barista::Providers::Provider do
  subject(:provider) do
    described_class.new(
      name: "aws",
      services: [Barista::Providers::Service.new(provider_name: "aws", name: "s3", url: "https://example.com")]
    )
  end

  it "exposes its name" do
    expect(provider.name).to eq("aws")
  end

  it "exposes its services" do
    expect(provider.services.size).to eq(1)
  end
end
