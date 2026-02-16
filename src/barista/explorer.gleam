/// The main agent loop.
///
/// Discover → Extract → Structure → Repeat.
///
/// Each step will eventually call out to scrape a provider's
/// SLA page and extract structured guarantees via an LLM.
import barista/guarantee

pub type Discovery {
  Discovery(
    provider: String,
    service: String,
    summary: String,
    guarantee: guarantee.Guarantee,
  )
}

/// Run the explorer agent against a list of provider SLA URLs.
/// Returns whatever guarantees it discovers.
pub fn run() -> List(Discovery) {
  discover()
}

/// The core discovery loop — not yet implemented.
fn discover() -> List(Discovery) {
  []
}
