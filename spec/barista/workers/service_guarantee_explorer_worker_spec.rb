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

    before do
      Barista.configure(
        Barista::Configuration.new(
          output_dir: dir,
          providers: [
            Barista::Providers::Provider.new(
              name: "aws",
              services: [
                Barista::Providers::Service.new(provider_name: "aws", name: "s3", url: "https://aws.example.com/s3/sla")
              ]
            )
          ]
        )
      )

      stub_request(:get, "https://aws.example.com/s3/sla")
        .to_return(status: 200, body: "<html>S3 SLA</html>")
    end

    after { FileUtils.rm_rf(dir) }

    it "writes a .caffeine file" do
      described_class.new.perform("aws/s3")
      expect(File.exist?(File.join(dir, "expectations", "s3.caffeine"))).to be(true)
    end

    it "generates content containing the service name" do
      described_class.new.perform("aws/s3")
      content = File.read(File.join(dir, "expectations", "s3.caffeine"))
      expect(content).to include('expectation "s3"')
    end
  end
end
