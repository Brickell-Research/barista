import gleam/option.{type Option}

pub type Metric {
  Availability
  Latency
  Throughput
  SuccessRate
}

pub type Comparator {
  GreaterThan
  GreaterThanOrEqualTo
  LessThan
  LessThanOrEqualTo
}

pub type Window {
  Monthly
  Annual
  Rolling(days: Int)
}

pub type Guarantee {
  Guarantee(
    provider: String,
    service: String,
    tier: Option(String),
    metric: Metric,
    threshold: Float,
    comparator: Comparator,
    window: Window,
    source_url: String,
  )
}
