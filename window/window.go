// Package window provides utilities for computing cron execution windows:
// the earliest and latest possible run times within a bounded time range.
package window

import (
	"fmt"
	"time"

	"github.com/cronparse/forecast"
	"github.com/cronparse/parser"
)

// Result holds the computed window boundaries for a cron expression.
type Result struct {
	Expression string
	From       time.Time
	To         time.Time
	First      *time.Time
	Last       *time.Time
	Count      int
	Err        error
}

// Compute returns the first and last run times of expr within [from, to],
// along with the total number of runs in that window.
func Compute(expr string, from, to time.Time) Result {
	res := Result{Expression: expr, From: from, To: to}

	if _, err := parser.Parse(expr); err != nil {
		res.Err = fmt.Errorf("invalid expression %q: %w", expr, err)
		return res
	}

	if !to.After(from) {
		res.Err = fmt.Errorf("to must be after from")
		return res
	}

	// Walk forward from `from` collecting all runs up to `to`.
	cursor := from
	for {
		runs := forecast.NextRuns(expr, cursor, 1)
		if len(runs) == 0 {
			break
		}
		t := runs[0]
		if t.After(to) {
			break
		}
		if res.First == nil {
			copy := t
			res.First = &copy
		}
		last := t
		res.Last = &last
		res.Count++
		cursor = t.Add(time.Second)
	}

	return res
}

// ComputeMany returns a Result for each expression in exprs.
func ComputeMany(exprs []string, from, to time.Time) []Result {
	results := make([]Result, len(exprs))
	for i, expr := range exprs {
		results[i] = Compute(expr, from, to)
	}
	return results
}
