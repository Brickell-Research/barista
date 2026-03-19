# frozen_string_literal: true

require "barista"

RSpec.describe Barista::Guarantee do
  subject(:guarantee) { described_class.new(name: "monthly_uptime", threshold: 99.9, window_days: 30) }

  it "exposes its attributes" do
    expect(guarantee).to have_attributes(name: "monthly_uptime", threshold: 99.9, window_days: 30)
  end

  it "is frozen" do
    expect(guarantee).to be_frozen
  end
end

RSpec.describe Barista::Intermediate do
  subject(:intermediate) do
    described_class.new(
      service_name: "s3",
      provider_name: "aws",
      source_url: "https://aws.amazon.com/s3/sla/",
      guarantees: [Barista::Guarantee.new(name: "monthly_uptime", threshold: 99.9, window_days: 30)]
    )
  end

  it "exposes its attributes" do
    expect(intermediate).to have_attributes(
      service_name: "s3",
      provider_name: "aws",
      source_url: "https://aws.amazon.com/s3/sla/"
    )
  end

  it "holds guarantees" do
    expect(intermediate.guarantees.size).to eq(1)
    expect(intermediate.guarantees.first).to have_attributes(name: "monthly_uptime", threshold: 99.9)
  end

  it "is frozen" do
    expect(intermediate).to be_frozen
  end
end
