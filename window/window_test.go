package window_test

import (
	"testing"
	"time"

	"github.com/cronparse/window"
)

func fixedFrom() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}

func TestCompute_EveryMinuteOneHour(t *testing.T) {
	from := fixedFrom()
	to := from.Add(time.Hour)
	res := window.Compute("* * * * *", from, to)

	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if res.Count != 60 {
		t.Errorf("expected 60 runs, got %d", res.Count)
	}
	if res.First == nil || res.Last == nil {
		t.Fatal("expected First and Last to be set")
	}
	if !res.First.Equal(from.Add(time.Minute)) {
		t.Errorf("unexpected First: %v", res.First)
	}
}

func TestCompute_HourlyTwoDays(t *testing.T) {
	from := fixedFrom()
	to := from.Add(48 * time.Hour)
	res := window.Compute("0 * * * *", from, to)

	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if res.Count != 48 {
		t.Errorf("expected 48 runs, got %d", res.Count)
	}
}

func TestCompute_InvalidExpression(t *testing.T) {
	from := fixedFrom()
	to := from.Add(time.Hour)
	res := window.Compute("invalid", from, to)

	if res.Err == nil {
		t.Error("expected error for invalid expression")
	}
	if res.Count != 0 {
		t.Errorf("expected 0 runs, got %d", res.Count)
	}
}

func TestCompute_ToBeforeFrom(t *testing.T) {
	from := fixedFrom()
	to := from.Add(-time.Hour)
	res := window.Compute("* * * * *", from, to)

	if res.Err == nil {
		t.Error("expected error when to is before from")
	}
}

func TestCompute_NoRunsInWindow(t *testing.T) {
	// Expression fires only on Jan 2; window is Jan 1 only.
	from := fixedFrom()
	to := from.Add(23*time.Hour + 59*time.Minute)
	res := window.Compute("0 12 2 1 *", from, to)

	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if res.Count != 0 {
		t.Errorf("expected 0 runs, got %d", res.Count)
	}
	if res.First != nil || res.Last != nil {
		t.Error("expected First and Last to be nil")
	}
}

func TestComputeMany_MixedExpressions(t *testing.T) {
	from := fixedFrom()
	to := from.Add(time.Hour)
	exprs := []string{"* * * * *", "0 * * * *", "bad"}
	results := window.ComputeMany(exprs, from, to)

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Count != 60 {
		t.Errorf("expected 60 for every-minute, got %d", results[0].Count)
	}
	if results[1].Count != 1 {
		t.Errorf("expected 1 for hourly, got %d", results[1].Count)
	}
	if results[2].Err == nil {
		t.Error("expected error for invalid expression")
	}
}
