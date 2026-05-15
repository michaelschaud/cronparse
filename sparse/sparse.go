// Package sparse identifies time windows where a cron expression fires
// infrequently (sparsely) relative to a given minimum gap threshold.
package sparse

import (
	"errors"
	"time"

	"github.com/cronparse/forecast"
)

// Gap represents a quiet period between two consecutive cron fires.
type Gap struct {
	From     time.Time
	To       time.Time
	Duration time.Duration
}

// Result holds the analysis output for a single expression.
type Result struct {
	Expression string
	Gaps       []Gap
	Longest    time.Duration
	Err        error
}

// Find returns all gaps between consecutive runs of expr within [from, to]
// that are strictly longer than minGap. At least two runs must exist in the
// window for any gaps to be detected.
func Find(expr string, from, to time.Time, minGap time.Duration) Result {
	res := Result{Expression: expr}

	if !to.After(from) {
		res.Err = errors.New("to must be after from")
		return res
	}

	runs, err := forecast.NextRuns(expr, from, 1440) // up to 1 day's worth of minutes
	if err != nil {
		res.Err = err
		return res
	}

	// Filter runs within [from, to]
	var filtered []time.Time
	for _, r := range runs {
		if !r.After(to) {
			filtered = append(filtered, r)
		}
	}

	if len(filtered) < 2 {
		return res
	}

	for i := 1; i < len(filtered); i++ {
		d := filtered[i].Sub(filtered[i-1])
		if d > minGap {
			g := Gap{
				From:     filtered[i-1],
				To:       filtered[i],
				Duration: d,
			}
			res.Gaps = append(res.Gaps, g)
			if d > res.Longest {
				res.Longest = d
			}
		}
	}

	return res
}
