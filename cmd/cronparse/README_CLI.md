# cronparse CLI

A command-line tool for parsing and forecasting cron expressions.

## Installation

```bash
go install github.com/example/cronparse/cmd/cronparse@latest
```

Or build from source:

```bash
git clone https://github.com/example/cronparse
cd cronparse
go build -o cronparse ./cmd/cronparse
```

## Usage

```
cronparse [flags] "<cron expression>"

Flags:
  -n int
        number of next run times to display (default 5)
  -from string
        start time for forecasting in RFC3339 format (default: current time)
```

## Examples

### Every minute
```bash
$ cronparse "* * * * *"
Expression : * * * * *
Description: every minute

Next 5 runs (from 2024-01-15 10:30:00):
  1. 2024-01-15 10:31:00 (Mon)
  2. 2024-01-15 10:32:00 (Mon)
  3. 2024-01-15 10:33:00 (Mon)
  4. 2024-01-15 10:34:00 (Mon)
  5. 2024-01-15 10:35:00 (Mon)
```

### Weekdays at 9am
```bash
$ cronparse -n 3 "0 9 * * 1-5"
Expression : 0 9 * * 1-5
Description: at 09:00, Monday through Friday

Next 3 runs (from 2024-01-15 10:30:00):
  1. 2024-01-16 09:00:00 (Tue)
  2. 2024-01-17 09:00:00 (Wed)
  3. 2024-01-18 09:00:00 (Thu)
```

### Forecast from a specific time
```bash
$ cronparse -n 3 -from 2024-06-01T00:00:00Z "0 0 * * 0"
Expression : 0 0 * * 0
Description: at 00:00, only on Sunday

Next 3 runs (from 2024-06-01 00:00:00):
  1. 2024-06-02 00:00:00 (Sun)
  2. 2024-06-09 00:00:00 (Sun)
  3. 2024-06-16 00:00:00 (Sun)
```

## Cron Expression Format

```
┌───────────── minute (0–59)
│ ┌───────────── hour (0–23)
│ │ ┌───────────── day of month (1–31)
│ │ │ ┌───────────── month (1–12)
│ │ │ │ ┌───────────── day of week (0–6, Sunday=0)
│ │ │ │ │
* * * * *
```

Supported syntax: `*`, numbers, ranges (`1-5`), steps (`*/15`), lists (`1,3,5`).
