package forecast

import (
	"fmt"
	"time"

	"github.com/cronparse/parser"
)

// NextRuns returns the next n scheduled run times after the given start time.
func NextRuns(expr string, start time.Time, n int) ([]time.Time, error) {
	if n <= 0 {
		return nil, fmt.Errorf("n must be greater than 0")
	}

	schedule, err := parser.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("invalid cron expression: %w", err)
	}

	results := make([]time.Time, 0, n)
	t := start.Truncate(time.Minute).Add(time.Minute)

	for len(results) < n {
		if matches(schedule, t) {
			results = append(results, t)
		}
		t = t.Add(time.Minute)
		// Guard against infinite loop for pathological expressions
		if t.After(start.Add(4 * 365 * 24 * time.Hour)) {
			return results, fmt.Errorf("could not find %d runs within 4 years", n)
		}
	}

	return results, nil
}

// NextRun returns the single next scheduled run time after the given start time.
func NextRun(expr string, start time.Time) (time.Time, error) {
	runs, err := NextRuns(expr, start, 1)
	if err != nil {
		return time.Time{}, err
	}
	return runs[0], nil
}

// matches reports whether t satisfies all fields of the parsed schedule.
func matches(s *parser.Schedule, t time.Time) bool {
	return fieldMatches(s.Minute, t.Minute()) &&
		fieldMatches(s.Hour, t.Hour()) &&
		fieldMatches(s.DayOfMonth, t.Day()) &&
		fieldMatches(s.Month, int(t.Month())) &&
		fieldMatches(s.DayOfWeek, int(t.Weekday()))
}

// fieldMatches checks if value is present in the allowed set.
func fieldMatches(allowed []int, value int) bool {
	for _, v := range allowed {
		if v == value {
			return true
		}
	}
	return false
}
