// Package slot finds available (non-firing) time slots between cron schedules
// within a given time window.
package slot

import (
	"errors"
	"time"

	"github.com/cronparse/forecast"
	"github.com/cronparse/parser"
)

// Gap represents a continuous period during which none of the given cron
// expressions fire.
type Gap struct {
	Start time.Time
	End   time.Time
}

// Duration returns the length of the gap.
func (g Gap) Duration() time.Duration {
	return g.End.Sub(g.Start)
}

// Find returns all gaps (quiet periods) within [from, to) where none of the
// provided cron expressions fire. Resolution controls the granularity of the
// scan (typically time.Minute).
func Find(expressions []string, from, to time.Time, resolution time.Duration) ([]Gap, error) {
	if to.Before(from) || to.Equal(from) {
		return nil, errors.New("slot: 'to' must be after 'from'")
	}
	if resolution <= 0 {
		return nil, errors.New("slot: resolution must be positive")
	}

	// Validate all expressions up front.
	for _, expr := range expressions {
		if _, err := parser.Parse(expr); err != nil {
			return nil, err
		}
	}

	var gaps []Gap
	var gapStart *time.Time

	for t := from; t.Before(to); t = t.Add(resolution) {
		fires := false
		for _, expr := range expressions {
			if firesAt(expr, t) {
				fires = true
				break
			}
		}

		if !fires {
			if gapStart == nil {
				copy := t
				gapStart = &copy
			}
		} else {
			if gapStart != nil {
				gaps = append(gaps, Gap{Start: *gapStart, End: t})
				gapStart = nil
			}
		}
	}

	// Close any open gap at the window boundary.
	if gapStart != nil {
		gaps = append(gaps, Gap{Start: *gapStart, End: to})
	}

	return gaps, nil
}

// firesAt reports whether expr fires at the exact minute represented by t.
// It does so by checking whether t is the next run at or immediately after
// t-1ns (i.e. the next run from just before t equals t).
func firesAt(expr string, t time.Time) bool {
	next, err := forecast.NextRun(expr, t.Add(-time.Nanosecond))
	if err != nil {
		return false
	}
	return next.Truncate(time.Minute).Equal(t.Truncate(time.Minute))
}
