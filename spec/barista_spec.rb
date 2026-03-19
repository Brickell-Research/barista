# frozen_string_literal: true

require "barista"

RSpec.describe Barista do
  it "is defined as a module" do
    expect(described_class).to be_a(Module)
  end
end
