# frozen_string_literal: true

require "barista"

RSpec.describe Barista::Translator do
  let(:guarantee) { Barista::Guarantee.new(name: "monthly_uptime_percentage", threshold: 99.9, window_days: 30) }
  let(:intermediate) do
    Barista::Intermediate.new(
      service_name: "s3",
      provider_name: "aws",
      source_url: "https://aws.amazon.com/s3/sla/",
      guarantees: [guarantee]
    )
  end

  describe ".translate" do
    subject(:output) { described_class.translate(intermediate) }

    it "opens with the Unmeasured Expectations header" do
      expect(output).to start_with("Unmeasured Expectations")
    end

    it "includes the guarantee name" do
      expect(output).to include('* "monthly_uptime_percentage"')
    end

    it "formats the threshold as a percentage" do
      expect(output).to include("threshold: 99.9%")
    end

    it "includes the window" do
      expect(output).to include("window_in_days: 30")
    end

    it "produces a valid Provides block" do
      expect(output).to include("    Provides {\n      threshold: 99.9%,\n      window_in_days: 30\n    }")
    end

    context "with multiple guarantees" do
      let(:intermediate) do
        Barista::Intermediate.new(
          service_name: "s3",
          provider_name: "aws",
          source_url: "https://aws.amazon.com/s3/sla/",
          guarantees: [
            Barista::Guarantee.new(name: "standard_availability", threshold: 99.9, window_days: 30),
            Barista::Guarantee.new(name: "reduced_redundancy", threshold: 99.0, window_days: 30)
          ]
        )
      end

      it "includes all guarantees" do
        expect(output).to include('* "standard_availability"')
        expect(output).to include('* "reduced_redundancy"')
      end

      it "formats whole-number thresholds without a decimal" do
        expect(output).to include("threshold: 99%")
      end
    end

    context "when there are no guarantees" do
      let(:intermediate) do
        Barista::Intermediate.new(
          service_name: "s3",
          provider_name: "aws",
          source_url: "https://aws.amazon.com/s3/sla/",
          guarantees: []
        )
      end

      it "does not include the Unmeasured Expectations header" do
        expect(output).not_to include("Unmeasured Expectations")
      end

      it "includes a comment with the service key" do
        expect(output).to include("# No guarantees found for aws/s3")
      end

      it "includes the source URL" do
        expect(output).to include("# Source: https://aws.amazon.com/s3/sla/")
      end
    end
  end
end
