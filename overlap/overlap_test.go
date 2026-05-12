package overlap_test

import (
	"testing"
	"time"

	"github.com/yourorg/cronparse/overlap"
)

var base = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func window(minutes int) time.Time {
	return base.Add(time.Duration(minutes) * time.Minute)
}

func TestFind_EveryMinuteOverlapsEveryMinute(t *testing.T) {
	r, err := overlap.Find("* * * * *", "* * * * *", base, window(5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Count != 5 {
		t.Errorf("expected 5 overlaps, got %d", r.Count)
	}
}

func TestFind_DisjointExpressions(t *testing.T) {
	// A fires at minute 0, B fires at minute 30
	r, err := overlap.Find("0 * * * *", "30 * * * *", base, window(60))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Count != 0 {
		t.Errorf("expected 0 overlaps, got %d", r.Count)
	}
}

func TestFind_PartialOverlap(t *testing.T) {
	// Both fire at minute 0 of every hour
	r, err := overlap.Find("0 * * * *", "0 */2 * * *", base, base.Add(4*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Count == 0 {
		t.Error("expected at least one overlap")
	}
}

func TestFind_InvalidExpressionA(t *testing.T) {
	_, err := overlap.Find("bad expr", "* * * * *", base, window(60))
	if err == nil {
		t.Error("expected error for invalid expression A")
	}
}

func TestFind_InvalidExpressionB(t *testing.T) {
	_, err := overlap.Find("* * * * *", "bad expr", base, window(60))
	if err == nil {
		t.Error("expected error for invalid expression B")
	}
}

func TestFind_InvalidWindow(t *testing.T) {
	_, err := overlap.Find("* * * * *", "* * * * *", window(60), base)
	if err == nil {
		t.Error("expected error when to is before from")
	}
}
