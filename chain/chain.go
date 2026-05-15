// Package chain provides utilities for chaining multiple cron expressions
// and computing a unified ordered sequence of next-run times across all of them.
package chain

import (
	"fmt"
	"sort"
	"time"

	"github.com/yourorg/cronparse/forecast"
	"github.com/yourorg/cronparse/parser"
)

// Entry holds a labeled cron expression and its resolved next-run times.
type Entry struct {
	Label      string
	Expression string
	Runs       []time.Time
	Err        error
}

// Result is the output of a Chain operation.
type Result struct {
	Entries  []Entry
	Merged   []MergedRun
}

// MergedRun associates a single scheduled time with the labels that fire at that moment.
type MergedRun struct {
	At     time.Time
	Labels []string
}

// Build resolves n next-run times for each labeled expression starting from
// from, then merges them into a chronologically sorted, deduplicated sequence.
func Build(from time.Time, n int, expressions map[string]string) Result {
	entries := make([]Entry, 0, len(expressions))

	// stable ordering by label
	labels := make([]string, 0, len(expressions))
	for l := range expressions {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	mergeMap := make(map[time.Time][]string)

	for _, label := range labels {
		expr := expressions[label]
		e := Entry{Label: label, Expression: expr}

		if _, err := parser.Parse(expr); err != nil {
			e.Err = fmt.Errorf("invalid expression %q: %w", expr, err)
			entries = append(entries, e)
			continue
		}

		runs, err := forecast.NextRuns(expr, from, n)
		if err != nil {
			e.Err = err
			entries = append(entries, e)
			continue
		}

		e.Runs = runs
		entries = append(entries, e)

		for _, r := range runs {
			mergeMap[r] = append(mergeMap[r], label)
		}
	}

	merged := make([]MergedRun, 0, len(mergeMap))
	for t, lbls := range mergeMap {
		merged = append(merged, MergedRun{At: t, Labels: lbls})
	}
	sort.Slice(merged, func(i, j int) bool {
		return merged[i].At.Before(merged[j].At)
	})

	return Result{Entries: entries, Merged: merged}
}
