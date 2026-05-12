// Package overlap provides utilities for detecting and reporting
// temporal overlaps between two cron expressions within a time window.
package overlap

import (
	"fmt"
	"time"

	"github.com/yourorg/cronparse/forecast"
	"github.com/yourorg/cronparse/parser"
)

// Result holds the overlap analysis between two cron expressions.
type Result struct {
	ExpressionA string
	ExpressionB string
	Overlaps    []time.Time
	Count       int
}

// Find returns all moments within [from, to) where both expressions
// would fire at the same minute boundary.
func Find(exprA, exprB string, from, to time.Time) (*Result, error) {
	if _, err := parser.Parse(exprA); err != nil {
		return nil, fmt.Errorf("expression A invalid: %w", err)
	}
	if _, err := parser.Parse(exprB); err != nil {
		return nil, fmt.Errorf("expression B invalid: %w", err)
	}
	if !to.After(from) {
		return nil, fmt.Errorf("to must be after from")
	}

	window := to.Sub(from)
	nA := int(window.Minutes()) + 1

	runsA, err := forecast.NextRuns(exprA, from, nA)
	if err != nil {
		return nil, fmt.Errorf("forecasting A: %w", err)
	}

	setA := make(map[time.Time]struct{}, len(runsA))
	for _, t := range runsA {
		if !t.Before(from) && t.Before(to) {
			setA[t.Truncate(time.Minute)] = struct{}{}
		}
	}

	runsB, err := forecast.NextRuns(exprB, from, nA)
	if err != nil {
		return nil, fmt.Errorf("forecasting B: %w", err)
	}

	var overlaps []time.Time
	for _, t := range runsB {
		key := t.Truncate(time.Minute)
		if !t.Before(from) && t.Before(to) {
			if _, ok := setA[key]; ok {
				overlaps = append(overlaps, key)
			}
		}
	}

	return &Result{
		ExpressionA: exprA,
		ExpressionB: exprB,
		Overlaps:    overlaps,
		Count:       len(overlaps),
	}, nil
}
