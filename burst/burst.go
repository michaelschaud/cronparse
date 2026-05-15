// Package burst identifies time windows where a cron expression fires
// more frequently than its average rate — useful for detecting clustering.
package burst

import (
	"fmt"
	"sort"
	"time"

	"github.com/cronparse/forecast"
)

// Window describes a contiguous burst of rapid firings.
type Window struct {
	Start     time.Time
	End       time.Time
	Firings   []time.Time
	GapMedian time.Duration
}

// Find returns burst windows within [from, to) for the given cron expression.
// A burst is a run of consecutive firings whose inter-firing gap is less than
// threshold. Minimum burst length is minFirings consecutive events.
func Find(expr string, from, to time.Time, threshold time.Duration, minFirings int) ([]Window, error) {
	if to.Before(from) || to.Equal(from) {
		return nil, fmt.Errorf("burst: to must be after from")
	}
	if minFirings < 2 {
		minFirings = 2
	}

	maxRuns := int(to.Sub(from)/time.Minute) + 1
	if maxRuns < 1 {
		maxRuns = 1
	}

	runs, err := forecast.NextRuns(expr, from, maxRuns)
	if err != nil {
		return nil, fmt.Errorf("burst: %w", err)
	}

	// filter to window
	var firings []time.Time
	for _, r := range runs {
		if !r.Before(from) && r.Before(to) {
			firings = append(firings, r)
		}
	}

	if len(firings) < minFirings {
		return nil, nil
	}

	var windows []Window
	start := 0
	for i := 1; i <= len(firings); i++ {
		inBurst := i < len(firings) && firings[i].Sub(firings[i-1]) < threshold
		if !inBurst {
			span := firings[start:i]
			if len(span) >= minFirings {
				windows = append(windows, Window{
					Start:     span[0],
					End:       span[len(span)-1],
					Firings:   append([]time.Time(nil), span...),
					GapMedian: medianGap(span),
				})
			}
			start = i
		}
	}

	return windows, nil
}

func medianGap(firings []time.Time) time.Duration {
	if len(firings) < 2 {
		return 0
	}
	gaps := make([]time.Duration, len(firings)-1)
	for i := 1; i < len(firings); i++ {
		gaps[i-1] = firings[i].Sub(firings[i-1])
	}
	sort.Slice(gaps, func(a, b int) bool { return gaps[a] < gaps[b] })
	return gaps[len(gaps)/2]
}
