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

    it "opens with the Expectations header" do
      expect(output).to start_with("Expectations")
    end

    it "includes the guarantee name" do
      expect(output).to include('"monthly_uptime_percentage"')
    end

    it "formats the guarantee as a Guarantees clause" do
      expect(output).to include("Guarantees 99.9% over 30d window")
    end

    it "indents the Guarantees clause under the name" do
      expect(output).to include("\"monthly_uptime_percentage\":\n  Guarantees 99.9% over 30d window")
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
        expect(output).to include('"standard_availability"')
        expect(output).to include('"reduced_redundancy"')
      end

      it "formats whole-number thresholds without a decimal" do
        expect(output).to include("Guarantees 99% over 30d window")
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

      it "does not include the Expectations header" do
        expect(output).not_to include("Expectations")
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
