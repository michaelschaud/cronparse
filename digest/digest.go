// Package digest provides a summary digest of a cron expression,
// including its frequency category, next run, and a human-readable description.
package digest

import (
	"fmt"
	"time"

	"github.com/cronparse/forecast"
	"github.com/cronparse/humanize"
	"github.com/cronparse/schedule"
	"github.com/cronparse/validate"
)

// Result holds the digest summary for a single cron expression.
type Result struct {
	Expression  string
	Valid       bool
	Error       string
	Description string
	Frequency   string
	NextRun     time.Time
	NextRunIn   time.Duration
}

// Of computes a digest for the given cron expression relative to from.
func Of(expr string, from time.Time) Result {
	if err := validate.Check(expr); err != nil {
		return Result{
			Expression: expr,
			Valid:      false,
			Error:      err.Error(),
		}
	}

	desc, err := humanize.Describe(expr)
	if err != nil {
		desc = ""
	}

	freq, err := schedule.Frequency(expr, from, from.Add(7*24*time.Hour), 500)
	freqLabel := ""
	if err == nil {
		freqLabel = classifyFrequency(freq)
	}

	next, err := forecast.NextRun(expr, from)
	if err != nil {
		return Result{
			Expression:  expr,
			Valid:       true,
			Description: desc,
			Frequency:   freqLabel,
		}
	}

	return Result{
		Expression:  expr,
		Valid:       true,
		Description: desc,
		Frequency:   freqLabel,
		NextRun:     next,
		NextRunIn:   next.Sub(from).Truncate(time.Second),
	}
}

// classifyFrequency returns a human label for a runs-per-week count.
func classifyFrequency(runsPerWeek int) string {
	switch {
	case runsPerWeek >= 10000:
		return "every minute"
	case runsPerWeek >= 1680:
		return "sub-hourly"
	case runsPerWeek >= 168:
		return "hourly"
	case runsPerWeek >= 24:
		return "daily"
	case runsPerWeek >= 7:
		return "weekly"
	case runsPerWeek >= 1:
		return "infrequent"
	default:
		return fmt.Sprintf("%d runs/week", runsPerWeek)
	}
}
