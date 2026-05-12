// Package closest finds the cron expression from a set that will next fire
// closest to (or furthest from) a given reference time.
package closest

import (
	"errors"
	"time"

	"github.com/your-org/cronparse/forecast"
)

// ErrNoExpressions is returned when the input slice is empty.
var ErrNoExpressions = errors.New("closest: no expressions provided")

// Result holds a cron expression together with its next scheduled run.
type Result struct {
	Expression string
	Next       time.Time
}

// Nearest returns the Result whose next run is closest to ref.
// All expressions are evaluated relative to ref.
func Nearest(expressions []string, ref time.Time) (Result, error) {
	if len(expressions) == 0 {
		return Result{}, ErrNoExpressions
	}

	best := Result{}
	for _, expr := range expressions {
		t, err := forecast.NextRun(expr, ref)
		if err != nil {
			return Result{}, err
		}
		if best.Expression == "" || t.Before(best.Next) {
			best = Result{Expression: expr, Next: t}
		}
	}
	return best, nil
}

// Farthest returns the Result whose next run is furthest from ref.
func Farthest(expressions []string, ref time.Time) (Result, error) {
	if len(expressions) == 0 {
		return Result{}, ErrNoExpressions
	}

	best := Result{}
	for _, expr := range expressions {
		t, err := forecast.NextRun(expr, ref)
		if err != nil {
			return Result{}, err
		}
		if best.Expression == "" || t.After(best.Next) {
			best = Result{Expression: expr, Next: t}
		}
	}
	return best, nil
}

// All returns Results for every expression sorted by next-run time (ascending).
func All(expressions []string, ref time.Time) ([]Result, error) {
	if len(expressions) == 0 {
		return nil, ErrNoExpressions
	}

	results := make([]Result, 0, len(expressions))
	for _, expr := range expressions {
		t, err := forecast.NextRun(expr, ref)
		if err != nil {
			return nil, err
		}
		results = append(results, Result{Expression: expr, Next: t})
	}

	// insertion sort — typical input is small
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Next.Before(results[j-1].Next); j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
	return results, nil
}
