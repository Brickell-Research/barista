# typed: false
# frozen_string_literal: true

require "barista"
require "sidekiq/testing"

RSpec.describe Barista::Workers::ServiceGuaranteeExplorerWorker do
  before do
    Sidekiq::Testing.fake!
    described_class.clear
  end

  it "uses the exploration queue" do
    expect(described_class.sidekiq_options["queue"]).to eq("exploration")
  end

  context "when called without a service name" do
    it "enqueues a job for each discovered service", :aggregate_failures do
      worker = described_class.new
      allow(worker).to receive(:discovered_services).and_return(%w[stripe twilio])

      worker.perform

      expect(described_class.jobs.size).to eq(2)
      expect(described_class.jobs.map { |j| j["args"] }).to contain_exactly(["stripe"], ["twilio"])
    end
  end

  context "when called with a service name" do
    it "explores that specific service" do
      worker = described_class.new
      allow(worker).to receive(:explore_service)

      worker.perform("stripe")

      expect(worker).to have_received(:explore_service).with("stripe")
    end
  end
end
