// Package history provides utilities for determining how many times
// a cron expression would have fired within a given time window.
package history

import (
	"fmt"
	"time"

	"github.com/cronparse/forecast"
	"github.com/cronparse/parser"
)

// Occurrence represents a single past firing of a cron expression.
type Occurrence struct {
	Time time.Time
}

// Result holds the full history of occurrences within a time window.
type Result struct {
	Expression string
	From       time.Time
	To         time.Time
	Occurrences []Occurrence
}

// Count returns the number of times the given cron expression would have
// fired between from (inclusive) and to (exclusive).
func Count(expr string, from, to time.Time) (int, error) {
	r, err := Between(expr, from, to)
	if err != nil {
		return 0, err
	}
	return len(r.Occurrences), nil
}

// Between returns all occurrences of the cron expression firing between
// from (inclusive) and to (exclusive).
func Between(expr string, from, to time.Time) (*Result, error) {
	if _, err := parser.Parse(expr); err != nil {
		return nil, fmt.Errorf("history: invalid expression %q: %w", expr, err)
	}
	if !to.After(from) {
		return nil, fmt.Errorf("history: 'to' must be after 'from'")
	}

	result := &Result{
		Expression: expr,
		From:       from,
		To:         to,
	}

	cursor := from.Truncate(time.Minute)
	for {
		runs, err := forecast.NextRuns(expr, cursor, 1)
		if err != nil || len(runs) == 0 {
			break
		}
		next := runs[0]
		if !next.Before(to) {
			break
		}
		result.Occurrences = append(result.Occurrences, Occurrence{Time: next})
		cursor = next.Add(time.Minute)
	}

	return result, nil
}
