// Package closest selects the cron expression(s) from a candidate set whose
// next scheduled run is nearest to (or farthest from) a reference time.
//
// # Overview
//
// Given a slice of cron expressions and a reference time, the package computes
// the next firing time for each expression and ranks them accordingly.
//
// # Functions
//
//   - Nearest  – returns the expression that fires soonest after ref.
//   - Farthest – returns the expression whose next run is furthest from ref.
//   - All      – returns all results sorted by next-run time (ascending).
//
// # Formatting
//
// FormatText and FormatMarkdown render a []Result as a plain-text table or a
// GitHub-flavoured Markdown table respectively.
//
// # Errors
//
// ErrNoExpressions is returned when the input slice is empty. Any invalid cron
// expression causes an immediate error from the underlying forecast package.
package closest
