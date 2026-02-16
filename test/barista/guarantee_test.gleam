import barista/guarantee.{Availability, GreaterThan, Guarantee, Monthly}
import gleam/option.{Some}
import gleeunit/should

pub fn guarantee_test() {
  let g =
    Guarantee(
      provider: "aws",
      service: "s3",
      tier: Some("standard"),
      metric: Availability,
      threshold: 99.9,
      comparator: GreaterThan,
      window: Monthly,
      source_url: "https://aws.amazon.com/s3/sla/",
    )

  g.provider |> should.equal("aws")
  g.threshold |> should.equal(99.9)
}
