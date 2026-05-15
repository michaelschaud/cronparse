// Package heatmap builds a frequency heatmap of cron expression fire times
// over a given time window, bucketed by hour-of-day and day-of-week.
package heatmap

import (
	"fmt"
	"time"

	"github.com/cronparse/forecast"
)

// Cell represents the number of times a cron expression fires within a
// particular hour-of-day (0–23) and day-of-week (0=Sunday … 6=Saturday) bucket.
type Cell struct {
	DayOfWeek int // 0 = Sunday, 6 = Saturday
	Hour      int // 0–23
	Count     int
}

// Map holds the full heatmap result for a single expression.
type Map struct {
	Expression string
	From       time.Time
	To         time.Time
	// Cells contains only buckets with Count > 0, sorted by DayOfWeek then Hour.
	Cells []Cell
	// Grid is a [7][24] matrix: Grid[dayOfWeek][hour] = count.
	Grid [7][24]int
}

// Build computes a heatmap for expr over the half-open window [from, to).
// It returns an error if expr is invalid or to is not after from.
func Build(expr string, from, to time.Time) (*Map, error) {
	if !to.After(from) {
		return nil, fmt.Errorf("heatmap: to must be after from")
	}

	// Collect all fire times in the window.
	runs, err := forecast.NextRuns(expr, from, int(to.Sub(from).Minutes())+1)
	if err != nil {
		return nil, fmt.Errorf("heatmap: %w", err)
	}

	m := &Map{
		Expression: expr,
		From:       from,
		To:         to,
	}

	for _, t := range runs {
		if t.Before(from) || !t.Before(to) {
			continue
		}
		dow := int(t.Weekday()) // Sunday=0
		h := t.Hour()
		m.Grid[dow][h]++
	}

	// Flatten non-zero cells.
	for dow := 0; dow < 7; dow++ {
		for h := 0; h < 24; h++ {
			if m.Grid[dow][h] > 0 {
				m.Cells = append(m.Cells, Cell{
					DayOfWeek: dow,
					Hour:      h,
					Count:     m.Grid[dow][h],
				})
			}
		}
	}

	return m, nil
}

// Peak returns the Cell with the highest Count. If there are no fire times,
// it returns a zero Cell.
func (m *Map) Peak() Cell {
	var best Cell
	for _, c := range m.Cells {
		if c.Count > best.Count {
			best = c
		}
	}
	return best
}
