// Package drift measures how much a cron expression "drifts" over time —
// i.e., the variance and standard deviation of intervals between consecutive
// scheduled runs within a given window.
package drift

import (
	"fmt"
	"math"
	"time"

	"github.com/example/cronparse/forecast"
	"github.com/example/cronparse/parser"
)

// Result holds drift statistics for a cron expression.
type Result struct {
	Expression  string
	SampleSize  int
	MeanGap     time.Duration
	MinGap      time.Duration
	MaxGap      time.Duration
	StdDev      time.Duration
	VarianceMs2 float64 // variance in ms²
}

// Measure computes drift statistics for expr over [from, to).
// It returns an error if the expression is invalid or fewer than 2 runs exist.
func Measure(expr string, from, to time.Time) (Result, error) {
	if _, err := parser.Parse(expr); err != nil {
		return Result{}, fmt.Errorf("drift: invalid expression %q: %w", expr, err)
	}
	if !to.After(from) {
		return Result{}, fmt.Errorf("drift: to must be after from")
	}

	// Collect all runs in the window.
	runs := forecast.NextRuns(expr, from, int(to.Sub(from).Minutes())+2)
	var filtered []time.Time
	for _, r := range runs {
		if !r.Before(from) && r.Before(to) {
			filtered = append(filtered, r)
		}
	}

	if len(filtered) < 2 {
		return Result{}, fmt.Errorf("drift: fewer than 2 runs in window")
	}

	gaps := make([]float64, len(filtered)-1)
	minGap := time.Duration(math.MaxInt64)
	maxGap := time.Duration(0)
	var sumMs float64

	for i := 1; i < len(filtered); i++ {
		gap := filtered[i].Sub(filtered[i-1])
		ms := float64(gap.Milliseconds())
		gaps[i-1] = ms
		sumMs += ms
		if gap < minGap {
			minGap = gap
		}
		if gap > maxGap {
			maxGap = gap
		}
	}

	meanMs := sumMs / float64(len(gaps))
	var variance float64
	for _, g := range gaps {
		d := g - meanMs
		variance += d * d
	}
	variance /= float64(len(gaps))
	stdDevMs := math.Sqrt(variance)

	return Result{
		Expression:  expr,
		SampleSize:  len(gaps),
		MeanGap:     time.Duration(meanMs) * time.Millisecond,
		MinGap:      minGap,
		MaxGap:      maxGap,
		StdDev:      time.Duration(stdDevMs) * time.Millisecond,
		VarianceMs2: variance,
	}, nil
}
