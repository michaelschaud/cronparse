// Package pivot identifies a "pivot" time — the next run of a reference
// expression — and then ranks a set of other expressions by how close their
// own next run lands relative to that pivot.
package pivot

import (
	"fmt"
	"sort"
	"time"

	"github.com/yourorg/cronparse/forecast"
)

// Result holds one expression together with its distance from the pivot.
type Result struct {
	Expression string
	NextRun    time.Time
	Delta      time.Duration // absolute distance from pivot; negative means before
	Error      error
}

// Around computes the pivot time from referenceExpr (its immediate next run
// after `from`) and returns one Result per expression in exprs, sorted by
// absolute distance from the pivot ascending.
func Around(referenceExpr string, exprs []string, from time.Time) (pivot time.Time, results []Result, err error) {
	pivotRuns, pivotErr := forecast.NextRuns(referenceExpr, from, 1)
	if pivotErr != nil {
		return time.Time{}, nil, fmt.Errorf("pivot expression invalid: %w", pivotErr)
	}
	if len(pivotRuns) == 0 {
		return time.Time{}, nil, fmt.Errorf("pivot expression produced no runs")
	}
	pivot = pivotRuns[0]

	results = make([]Result, 0, len(exprs))
	for _, expr := range exprs {
		runs, runErr := forecast.NextRuns(expr, from, 1)
		if runErr != nil {
			results = append(results, Result{Expression: expr, Error: runErr})
			continue
		}
		var next time.Time
		if len(runs) > 0 {
			next = runs[0]
		}
		delta := next.Sub(pivot)
		results = append(results, Result{
			Expression: expr,
			NextRun:    next,
			Delta:      delta,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Error != nil {
			return false
		}
		if results[j].Error != nil {
			return true
		}
		di := results[i].Delta
		if di < 0 {
			di = -di
		}
		dj := results[j].Delta
		if dj < 0 {
			dj = -dj
		}
		return di < dj
	})

	return pivot, results, nil
}
