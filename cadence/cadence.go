// Package cadence analyzes the regularity of a cron expression over a window
// of time, reporting whether runs occur at consistent intervals.
package cadence

import (
	"errors"
	"time"

	"github.com/cronparse/forecast"
)

// Result holds the cadence analysis for a single cron expression.
type Result struct {
	Expression string
	Regular    bool
	MinGap     time.Duration
	MaxGap     time.Duration
	MeanGap    time.Duration
	Jitter     time.Duration // MaxGap - MinGap
	SampleSize int
}

// Analyze computes the cadence of expr over the window [from, to].
// It returns an error if the expression is invalid or fewer than 2 runs
// occur in the window.
func Analyze(expr string, from, to time.Time, threshold time.Duration) (Result, error) {
	if to.Before(from) || to.Equal(from) {
		return Result{}, errors.New("cadence: 'to' must be after 'from'")
	}

	runs, err := forecast.NextRuns(expr, from, 1000)
	if err != nil {
		return Result{}, err
	}

	// filter to window
	var filtered []time.Time
	for _, r := range runs {
		if !r.After(to) {
			filtered = append(filtered, r)
		}
	}

	if len(filtered) < 2 {
		return Result{}, errors.New("cadence: fewer than 2 runs in the specified window")
	}

	gaps := make([]time.Duration, len(filtered)-1)
	for i := 1; i < len(filtered); i++ {
		gaps[i-1] = filtered[i].Sub(filtered[i-1])
	}

	minG, maxG, sumG := gaps[0], gaps[0], time.Duration(0)
	for _, g := range gaps {
		if g < minG {
			minG = g
		}
		if g > maxG {
			maxG = g
		}
		sumG += g
	}
	mean := sumG / time.Duration(len(gaps))
	jitter := maxG - minG

	return Result{
		Expression: expr,
		Regular:    jitter <= threshold,
		MinGap:     minG,
		MaxGap:     maxG,
		MeanGap:    mean,
		Jitter:     jitter,
		SampleSize: len(filtered),
	}, nil
}
