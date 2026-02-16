import barista/explorer
import gleam/list
import gleeunit/should

pub fn run_returns_empty_test() {
  explorer.run() |> list.length |> should.equal(0)
}
