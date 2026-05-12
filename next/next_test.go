package next_test

import (
	"testing"
	"time"

	"github.com/yourorg/cronparse/next"
)

var refTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func TestForExpressions_ValidExpressions(t *testing.T) {
	exprs := []string{"* * * * *", "0 * * * *"}
	results := next.ForExpressions(exprs, refTime, 3)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error for %q: %v", r.Expression, r.Err)
		}
		if len(r.Runs) != 3 {
			t.Errorf("expected 3 runs for %q, got %d", r.Expression, len(r.Runs))
		}
	}
}

func TestForExpressions_InvalidExpression(t *testing.T) {
	results := next.ForExpressions([]string{"invalid"}, refTime, 3)
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if results[0].Err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestForExpressions_RunsAreAfterFrom(t *testing.T) {
	results := next.ForExpressions([]string{"* * * * *"}, refTime, 5)
	for _, run := range results[0].Runs {
		if !run.After(refTime) {
			t.Errorf("run %v is not after reference time %v", run, refTime)
		}
	}
}

func TestMerged_DeduplicatesAndSorts(t *testing.T) {
	// Two identical expressions should deduplicate
	times, err := next.Merged([]string{"* * * * *", "* * * * *"}, refTime, 5, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 5 {
		t.Errorf("expected 5 deduplicated times, got %d", len(times))
	}
	for i := 1; i < len(times); i++ {
		if !times[i].After(times[i-1]) {
			t.Errorf("times not sorted at index %d", i)
		}
	}
}

func TestMerged_InvalidExpression(t *testing.T) {
	_, err := next.Merged([]string{"bad expr"}, refTime, 3, 10)
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestMerged_LimitApplied(t *testing.T) {
	times, err := next.Merged([]string{"* * * * *", "0 * * * *"}, refTime, 10, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) > 4 {
		t.Errorf("expected at most 4 times, got %d", len(times))
	}
}
