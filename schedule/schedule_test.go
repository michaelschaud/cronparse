package schedule_test

import (
	"testing"
	"time"

	"github.com/yourorg/cronparse/schedule"
)

var baseTime = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestCompare_OverlappingExpressions(t *testing.T) {
	diff, err := schedule.Compare("*/5 * * * *", "*/10 * * * *", baseTime, 12)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diff.NextA) != 12 {
		t.Errorf("expected 12 runs for A, got %d", len(diff.NextA))
	}
	if len(diff.NextB) != 12 {
		t.Errorf("expected 12 runs for B, got %d", len(diff.NextB))
	}
	// every 10-min is a subset of every 5-min, so OnlyInB should be empty
	if len(diff.OnlyInB) != 0 {
		t.Errorf("expected no runs only in B, got %d", len(diff.OnlyInB))
	}
	if len(diff.Common) == 0 {
		t.Error("expected some common runs")
	}
}

func TestCompare_DisjointExpressions(t *testing.T) {
	// minute 1 vs minute 2 — never overlap
	diff, err := schedule.Compare("1 * * * *", "2 * * * *", baseTime, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diff.Common) != 0 {
		t.Errorf("expected no common runs, got %d", len(diff.Common))
	}
	if len(diff.OnlyInA) != 5 {
		t.Errorf("expected 5 only-in-A runs, got %d", len(diff.OnlyInA))
	}
}

func TestCompare_InvalidExpressionA(t *testing.T) {
	_, err := schedule.Compare("bad expr", "* * * * *", baseTime, 5)
	if err == nil {
		t.Error("expected error for invalid expression A")
	}
}

func TestCompare_InvalidExpressionB(t *testing.T) {
	_, err := schedule.Compare("* * * * *", "99 * * * *", baseTime, 5)
	if err == nil {
		t.Error("expected error for invalid expression B")
	}
}
