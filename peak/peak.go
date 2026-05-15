// Package peak identifies the busiest time windows within a cron schedule
// by counting how many expressions fire within each sub-window of a given duration.
package peak

import (
	"fmt"
	"time"

	"github.com/cronparse/forecast"
)

// Window represents a time window with a count of how many cron fires occurred.
type Window struct {
	Start time.Time
	End   time.Time
	Count int
}

// Result holds the top busy windows found for a given expression.
type Result struct {
	Expression string
	Windows    []Window
	Error      error
}

// Find returns the top N busiest windows of the given duration within [from, to]
// for each provided cron expression. Windows are evaluated by sliding one slot
// at a time (step = windowSize / 2) to avoid missing boundary bursts.
func Find(expressions []string, from, to time.Time, windowSize time.Duration, topN int) []Result {
	results := make([]Result, 0, len(expressions))
	for _, expr := range expressions {
		results = append(results, findForExpr(expr, from, to, windowSize, topN))
	}
	return results
}

func findForExpr(expr string, from, to time.Time, windowSize time.Duration, topN int) Result {
	if to.Before(from) || to.Equal(from) {
		return Result{Expression: expr, Error: fmt.Errorf("to must be after from")}
	}
	if windowSize <= 0 {
		return Result{Expression: expr, Error: fmt.Errorf("windowSize must be positive")}
	}

	// Collect all fires in [from, to]
	allRuns, err := forecast.NextRuns(expr, from, int(to.Sub(from).Minutes())+2)
	if err != nil {
		return Result{Expression: expr, Error: err}
	}

	// Filter to [from, to]
	var fires []time.Time
	for _, r := range allRuns {
		if !r.Before(from) && !r.After(to) {
			fires = append(fires, r)
		}
	}

	step := windowSize / 2
	if step == 0 {
		step = time.Minute
	}

	var windows []Window
	for wStart := from; wStart.Before(to); wStart = wStart.Add(step) {
		wEnd := wStart.Add(windowSize)
		if wEnd.After(to) {
			wEnd = to
		}
		count := 0
		for _, f := range fires {
			if !f.Before(wStart) && f.Before(wEnd) {
				count++
			}
		}
		windows = append(windows, Window{Start: wStart, End: wEnd, Count: count})
	}

	// Sort descending by count (simple selection for small N)
	for i := 0; i < len(windows)-1; i++ {
		for j := i + 1; j < len(windows); j++ {
			if windows[j].Count > windows[i].Count {
				windows[i], windows[j] = windows[j], windows[i]
			}
		}
	}

	if topN > 0 && len(windows) > topN {
		windows = windows[:topN]
	}

	return Result{Expression: expr, Windows: windows}
}
