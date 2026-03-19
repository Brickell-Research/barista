# frozen_string_literal: true

require "barista"
require "tmpdir"

RSpec.describe Barista::CaffeineWriter do
  describe ".write" do
    let(:dir) { Dir.mktmpdir }
    let(:service) { Barista::Providers::Service.new(provider_name: "aws", name: "s3", url: "https://example.com") }

    before { Barista.configure(Barista::Configuration.new(providers: [], output_dir: dir)) }

    after { FileUtils.rm_rf(dir) }

    it "returns the expected file path" do
      path = described_class.write(service: service, content: "test")
      expect(path).to eq(File.join(dir, "expectations", "s3.caffeine"))
    end

    it "writes the content to disk" do
      path = described_class.write(service: service, content: 'expectation "s3" {}')
      expect(File.read(path)).to eq('expectation "s3" {}')
    end
  end
end
