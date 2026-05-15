// Package heatmap provides tools for visualising how frequently a cron
// expression fires across different hours of the day and days of the week.
//
// # Overview
//
// Build computes a [Map] by collecting all fire times for an expression within
// a caller-supplied time window and bucketing them into a [7][24] grid indexed
// by day-of-week (Sunday = 0) and hour-of-day (0–23).
//
// # Formatting
//
// Two rendering helpers are provided:
//
//   - [FormatText] produces a compact ASCII grid suitable for terminal output.
//   - [FormatMarkdown] produces a GitHub-flavoured Markdown table.
//
// # Example
//
//	from := time.Now().Truncate(24 * time.Hour)
//	to   := from.Add(7 * 24 * time.Hour)
//
//	m, err := heatmap.Build("0 9 * * 1-5", from, to)
//	if err != nil { log.Fatal(err) }
//
//	fmt.Println(heatmap.FormatText(m))
package heatmap
