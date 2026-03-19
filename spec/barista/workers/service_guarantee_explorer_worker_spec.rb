# frozen_string_literal: true

require "barista"
require "sidekiq/testing"
require "tmpdir"

RSpec.describe Barista::Workers::ServiceGuaranteeExplorerWorker do
  before do
    Sidekiq::Testing.fake!
    described_class.clear
  end

  it "uses the exploration queue" do
    expect(described_class.sidekiq_options["queue"]).to eq("exploration")
  end

  context "when called without a service key" do
    before do
      Barista.configure(
        Barista::Configuration.new(
          providers: [
            Barista::Providers::Provider.new(
              name: "stripe",
              services: [Barista::Providers::Service.new(provider_name: "stripe", name: "payments", url: "https://example.com")]
            ),
            Barista::Providers::Provider.new(
              name: "aws",
              services: [Barista::Providers::Service.new(provider_name: "aws", name: "s3", url: "https://example.com")]
            )
          ]
        )
      )
    end

    it "enqueues a job for each discovered service" do
      described_class.new.perform
      expect(described_class.jobs.size).to eq(2)
    end

    it "enqueues jobs with provider/service keys" do
      described_class.new.perform
      keys = described_class.jobs.map { |j| j["args"] }
      expect(keys).to contain_exactly(["stripe/payments"], ["aws/s3"])
    end
  end

  context "when called with a service key" do
    let(:dir) { Dir.mktmpdir }
    let(:service) do
      Barista::Providers::Service.new(provider_name: "aws", name: "s3", url: "https://aws.example.com/s3/sla")
    end
    let(:intermediate) do
      Barista::Intermediate.new(
        service_name: "s3",
        provider_name: "aws",
        source_url: "https://aws.example.com/s3/sla",
        guarantees: [Barista::Guarantee.new(name: "monthly_uptime_percentage", threshold: 99.9, window_days: 30)]
      )
    end

    before do
      Barista.configure(
        Barista::Configuration.new(
          output_dir: dir,
          providers: [
            Barista::Providers::Provider.new(name: "aws", services: [service])
          ]
        )
      )

      stub_request(:get, "https://aws.example.com/s3/sla")
        .to_return(status: 200, body: "<html>S3 SLA</html>")

      allow(Barista::Synthesizer).to receive(:synthesize).and_return(intermediate)
    end

    after { FileUtils.rm_rf(dir) }

    it "writes a .caffeine file" do
      described_class.new.perform("aws/s3")
      expect(File.exist?(File.join(dir, "expectations", "s3.caffeine"))).to be(true)
    end

    it "writes valid Caffeine content" do
      described_class.new.perform("aws/s3")
      content = File.read(File.join(dir, "expectations", "s3.caffeine"))
      expect(content).to include("Expectations")
      expect(content).to include('"monthly_uptime_percentage"')
      expect(content).to include("Guarantees 99.9% over 30d window")
    end

    context "when the service key is unknown" do
      it "does not raise and does not write a file" do
        described_class.new.perform("aws/unknown")
        expect(File.exist?(File.join(dir, "expectations", "unknown.caffeine"))).to be(false)
      end
    end
  end
end
