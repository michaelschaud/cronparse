# cronparse

A Go library and CLI for human-readable cron expression parsing and next-run forecasting.

---

## Installation

```bash
go get github.com/yourusername/cronparse
```

For the CLI:

```bash
go install github.com/yourusername/cronparse/cmd/cronparse@latest
```

---

## Usage

### Library

```go
package main

import (
    "fmt"
    "github.com/yourusername/cronparse"
)

func main() {
    expr, err := cronparse.Parse("*/15 9-17 * * 1-5")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(expr.Describe())
    // Output: Every 15 minutes, between 09:00 and 17:00, Monday through Friday

    next := expr.Next(5)
    for _, t := range next {
        fmt.Println(t)
    }
}
```

### CLI

```bash
# Describe a cron expression
cronparse describe "0 8 * * MON-FRI"
# Output: At 08:00, Monday through Friday

# Show the next N scheduled runs
cronparse next "*/30 * * * *" --count 5
```

---

## Features

- Human-readable descriptions of cron expressions
- Next-run forecasting for any number of future occurrences
- Supports standard 5-field cron syntax with named weekdays and months
- Simple CLI for quick inspection without writing code

---

## License

This project is licensed under the [MIT License](LICENSE).