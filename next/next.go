// Package next provides utilities for finding the next N scheduled
// run times for one or more cron expressions, with optional filtering
// and formatting support.
package next

import (
	"fmt"
	"sort"
	"time"

	"github.com/yourorg/cronparse/forecast"
	"github.com/yourorg/cronparse/parser"
)

// Result holds the next run times for a single cron expression.
type Result struct {
	Expression string
	Runs       []time.Time
	Err        error
}

// ForExpressions returns the next n run times for each expression,
// starting from the given reference time.
func ForExpressions(exprs []string, from time.Time, n int) []Result {
	results := make([]Result, 0, len(exprs))
	for _, expr := range exprs {
		r := Result{Expression: expr}
		_, err := parser.Parse(expr)
		if err != nil {
			r.Err = fmt.Errorf("invalid expression %q: %w", expr, err)
			results = append(results, r)
			continue
		}
		runs, err := forecast.NextRuns(expr, from, n)
		if err != nil {
			r.Err = err
		} else {
			r.Runs = runs
		}
		results = append(results, r)
	}
	return results
}

// Merged returns all next run times across all expressions, sorted
// chronologically and deduplicated, up to limit entries.
func Merged(exprs []string, from time.Time, n int, limit int) ([]time.Time, error) {
	seen := make(map[time.Time]struct{})
	var all []time.Time
	for _, r := range ForExpressions(exprs, from, n) {
		if r.Err != nil {
			return nil, r.Err
		}
		for _, t := range r.Runs {
			if _, ok := seen[t]; !ok {
				seen[t] = struct{}{}
				all = append(all, t)
			}
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Before(all[j]) })
	if limit > 0 && len(all) > limit {
		all = all[:limit]
	}
	return all, nil
}
