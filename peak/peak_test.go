package peak_test

import (
	"testing"
	"time"

	"github.com/cronparse/peak"
)

func fixedFrom() time.Time {
	return time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
}

func TestFind_EveryMinuteBusiestWindow(t *testing.T) {
	from := fixedFrom()
	to := from.Add(2 * time.Hour)
	results := peak.Find([]string{"* * * * *"}, from, to, 10*time.Minute, 3)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Error != nil {
		t.Fatalf("unexpected error: %v", r.Error)
	}
	if len(r.Windows) == 0 {
		t.Fatal("expected at least one window")
	}
	// Every-minute expression should yield ~10 fires per 10-min window
	if r.Windows[0].Count < 9 {
		t.Errorf("expected busiest window count >= 9, got %d", r.Windows[0].Count)
	}
}

func TestFind_TopNLimitsResults(t *testing.T) {
	from := fixedFrom()
	to := from.Add(3 * time.Hour)
	results := peak.Find([]string{"* * * * *"}, from, to, 15*time.Minute, 2)
	r := results[0]
	if r.Error != nil {
		t.Fatalf("unexpected error: %v", r.Error)
	}
	if len(r.Windows) > 2 {
		t.Errorf("expected at most 2 windows, got %d", len(r.Windows))
	}
}

func TestFind_InvalidExpression(t *testing.T) {
	from := fixedFrom()
	to := from.Add(1 * time.Hour)
	results := peak.Find([]string{"invalid expr"}, from, to, 10*time.Minute, 3)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Error == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestFind_ToBeforeFrom(t *testing.T) {
	from := fixedFrom()
	to := from.Add(-1 * time.Hour)
	results := peak.Find([]string{"* * * * *"}, from, to, 10*time.Minute, 3)
	if results[0].Error == nil {
		t.Error("expected error when to is before from")
	}
}

func TestFind_MultipleExpressions(t *testing.T) {
	from := fixedFrom()
	to := from.Add(2 * time.Hour)
	exprs := []string{"* * * * *", "0 * * * *"}
	results := peak.Find(exprs, from, to, 30*time.Minute, 3)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Error != nil {
			t.Errorf("unexpected error for %q: %v", r.Expression, r.Error)
		}
	}
	// every-minute should have more fires per window than hourly
	if results[0].Windows[0].Count <= results[1].Windows[0].Count {
		t.Errorf("expected every-minute to have more fires than hourly in busiest window")
	}
}

func TestFind_WindowsSortedDescending(t *testing.T) {
	from := fixedFrom()
	to := from.Add(4 * time.Hour)
	results := peak.Find([]string{"*/5 * * * *"}, from, to, 20*time.Minute, 5)
	r := results[0]
	if r.Error != nil {
		t.Fatalf("unexpected error: %v", r.Error)
	}
	for i := 1; i < len(r.Windows); i++ {
		if r.Windows[i].Count > r.Windows[i-1].Count {
			t.Errorf("windows not sorted descending at index %d", i)
		}
	}
}
