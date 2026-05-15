package pivot_test

import (
	"testing"
	"time"

	"github.com/yourorg/cronparse/pivot"
)

// fixedFrom returns a deterministic reference time: 2024-03-15 12:00:00 UTC.
func fixedFrom() time.Time {
	return time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)
}

func TestAround_InvalidPivotExpression(t *testing.T) {
	_, _, err := pivot.Around("invalid", []string{"* * * * *"}, fixedFrom())
	if err == nil {
		t.Fatal("expected error for invalid pivot expression, got nil")
	}
}

func TestAround_ReturnsPivotTime(t *testing.T) {
	// "* * * * *" fires every minute; pivot should be fixedFrom + 1 min.
	piv, _, err := pivot.Around("* * * * *", []string{}, fixedFrom())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := fixedFrom().Add(time.Minute)
	if !piv.Equal(expected) {
		t.Errorf("pivot = %v, want %v", piv, expected)
	}
}

func TestAround_ResultsSortedByAbsoluteDelta(t *testing.T) {
	// pivot fires at :01 past the hour (from 12:00 -> 12:01).
	// "0 * * * *" fires at :00 — that is 13:00, delta = +59 min.
	// "* * * * *" fires at :01, delta = 0.
	// "30 * * * *" fires at 12:30, delta = +29 min.
	exprs := []string{"0 * * * *", "* * * * *", "30 * * * *"}
	_, results, err := pivot.Around("* * * * *", exprs, fixedFrom())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	// Closest first: "* * * * *" (delta=0), then "30 * * * *", then "0 * * * *".
	for i := 1; i < len(results); i++ {
		prev := results[i-1].Delta
		if prev < 0 {
			prev = -prev
		}
		curr := results[i].Delta
		if curr < 0 {
			curr = -curr
		}
		if prev > curr {
			t.Errorf("results not sorted: index %d delta %v > index %d delta %v", i-1, prev, i, curr)
		}
	}
}

func TestAround_InvalidExprInSlice(t *testing.T) {
	_, results, err := pivot.Around("* * * * *", []string{"* * * * *", "bad expr"}, fixedFrom())
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	var errCount int
	for _, r := range results {
		if r.Error != nil {
			errCount++
		}
	}
	if errCount != 1 {
		t.Errorf("expected 1 errored result, got %d", errCount)
	}
}
