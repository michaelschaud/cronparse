// Package schedule provides utilities for building and comparing cron schedules,
// including diffing two expressions and summarizing schedule frequency.
package schedule

import (
	"fmt"
	"time"

	"github.com/yourorg/cronparse/forecast"
	"github.com/yourorg/cronparse/parser"
)

// Diff holds the comparison result between two cron expressions.
type Diff struct {
	ExprA     string
	ExprB     string
	NextA     []time.Time
	NextB     []time.Time
	OnlyInA   []time.Time
	OnlyInB   []time.Time
	Common    []time.Time
}

// Compare returns a Diff between two cron expressions over the next n occurrences
// starting from the given base time.
func Compare(exprA, exprB string, base time.Time, n int) (*Diff, error) {
	_, err := parser.Parse(exprA)
	if err != nil {
		return nil, fmt.Errorf("invalid expression A: %w", err)
	}
	_, err = parser.Parse(exprB)
	if err != nil {
		return nil, fmt.Errorf("invalid expression B: %w", err)
	}

	nextA, err := forecast.NextRuns(exprA, base, n)
	if err != nil {
		return nil, err
	}
	nextB, err := forecast.NextRuns(exprB, base, n)
	if err != nil {
		return nil, err
	}

	setA := toSet(nextA)
	setB := toSet(nextB)

	diff := &Diff{
		ExprA: exprA,
		ExprB: exprB,
		NextA: nextA,
		NextB: nextB,
	}

	for _, t := range nextA {
		if setB[t] {
			diff.Common = append(diff.Common, t)
		} else {
			diff.OnlyInA = append(diff.OnlyInA, t)
		}
	}
	for _, t := range nextB {
		if !setA[t] {
			diff.OnlyInB = append(diff.OnlyInB, t)
		}
	}

	return diff, nil
}

func toSet(times []time.Time) map[time.Time]bool {
	m := make(map[time.Time]bool, len(times))
	for _, t := range times {
		m[t] = true
	}
	return m
}
