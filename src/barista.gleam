import barista/explorer
import gleam/io

pub fn main() -> Nil {
  io.println("barista v0.1.0")
  io.println("")

  let results = explorer.run()

  case results {
    [] -> io.println("No guarantees discovered.")
    _ -> {
      io.println("Discovered guarantees:")
      print_guarantees(results)
    }
  }
}

fn print_guarantees(guarantees: List(explorer.Discovery)) -> Nil {
  case guarantees {
    [] -> Nil
    [first, ..rest] -> {
      io.println(
        "  " <> first.provider <> "/" <> first.service <> " " <> first.summary,
      )
      print_guarantees(rest)
    }
  }
}
