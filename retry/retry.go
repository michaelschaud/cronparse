// Package retry suggests alternative cron expressions when a given expression
// fires too infrequently or too frequently relative to a desired target cadence.
package retry

import (
	"fmt"
	"math"
	"time"

	"github.com/cronparse/forecast"
	"github.com/cronparse/parser"
)

// Suggestion holds an alternative cron expression and metadata about how
// closely its firing frequency matches the requested target interval.
type Suggestion struct {
	Expression  string
	Description string
	DeltaAbs    time.Duration // absolute difference from target interval
}

// candidates are pre-built expressions paired with a rough natural-language
// description; they cover the most common scheduling intervals.
var candidates = []struct {
	expr string
	desc string
}{
	{"* * * * *", "every minute"},
	{"*/2 * * * *", "every 2 minutes"},
	{"*/5 * * * *", "every 5 minutes"},
	{"*/10 * * * *", "every 10 minutes"},
	{"*/15 * * * *", "every 15 minutes"},
	{"*/30 * * * *", "every 30 minutes"},
	{"0 * * * *", "every hour"},
	{"0 */2 * * *", "every 2 hours"},
	{"0 */4 * * *", "every 4 hours"},
	{"0 */6 * * *", "every 6 hours"},
	{"0 */12 * * *", "every 12 hours"},
	{"0 0 * * *", "once a day at midnight"},
	{"0 9 * * *", "once a day at 09:00"},
	{"0 0 * * 1", "once a week on Monday"},
	{"0 0 1 * *", "once a month on the 1st"},
}

// Suggest returns up to n candidate expressions whose average firing interval
// is closest to target, excluding expr itself. The sample window used to
// estimate intervals is 7 days starting from now.
func Suggest(expr string, target time.Duration, n int) ([]Suggestion, error) {
	if _, err := parser.Parse(expr); err != nil {
		return nil, fmt.Errorf("retry: invalid expression %q: %w", expr, err)
	}
	if target <= 0 {
		return nil, fmt.Errorf("retry: target duration must be positive")
	}
	if n <= 0 {
		n = 3
	}

	from := time.Now().UTC().Truncate(time.Minute)
	window := 7 * 24 * time.Hour

	var results []Suggestion
	for _, c := range candidates {
		if c.expr == expr {
			continue
		}
		avg, err := avgInterval(c.expr, from, window)
		if err != nil || avg <= 0 {
			continue
		}
		delta := time.Duration(math.Abs(float64(avg - target)))
		results = append(results, Suggestion{
			Expression:  c.expr,
			Description: c.desc,
			DeltaAbs:    delta,
		})
	}

	// sort ascending by delta
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].DeltaAbs < results[j-1].DeltaAbs; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}

	if n < len(results) {
		results = results[:n]
	}
	return results, nil
}

// avgInterval estimates the mean firing interval for expr over the given window.
func avgInterval(expr string, from time.Time, window time.Duration) (time.Duration, error) {
	runs, err := forecast.NextRuns(expr, from, 50)
	if err != nil {
		return 0, err
	}
	to := from.Add(window)
	var filtered []time.Time
	for _, r := range runs {
		if !r.After(to) {
			filtered = append(filtered, r)
		}
	}
	if len(filtered) < 2 {
		return 0, nil
	}
	total := filtered[len(filtered)-1].Sub(filtered[0])
	return total / time.Duration(len(filtered)-1), nil
}
