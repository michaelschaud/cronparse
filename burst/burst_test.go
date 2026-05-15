package burst_test

import (
	"testing"
	"time"

	"github.com/cronparse/burst"
)

var fixedFrom = time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)

func TestFind_EveryMinuteBurstsWithSmallThreshold(t *testing.T) {
	to := fixedFrom.Add(10 * time.Minute)
	// every-minute expression with a 2-minute threshold — all consecutive
	windows, err := burst.Find("* * * * *", fixedFrom, to, 2*time.Minute, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) == 0 {
		t.Fatal("expected at least one burst window")
	}
	for _, w := range windows {
		if len(w.Firings) < 3 {
			t.Errorf("burst window has fewer than minFirings: %d", len(w.Firings))
		}
		if w.End.Before(w.Start) {
			t.Errorf("window end before start")
		}
	}
}

func TestFind_HourlyNoBurst(t *testing.T) {
	to := fixedFrom.Add(6 * time.Hour)
	// hourly fires once per hour; gap = 60 min; threshold = 30 min => no burst
	windows, err := burst.Find("0 * * * *", fixedFrom, to, 30*time.Minute, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) != 0 {
		t.Errorf("expected no burst windows, got %d", len(windows))
	}
}

func TestFind_InvalidExpression(t *testing.T) {
	to := fixedFrom.Add(time.Hour)
	_, err := burst.Find("invalid", fixedFrom, to, time.Minute, 2)
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestFind_ToBeforeFrom(t *testing.T) {
	_, err := burst.Find("* * * * *", fixedFrom, fixedFrom.Add(-time.Minute), time.Minute, 2)
	if err == nil {
		t.Fatal("expected error when to is before from")
	}
}

func TestFind_GapMedianIsPositive(t *testing.T) {
	to := fixedFrom.Add(20 * time.Minute)
	windows, err := burst.Find("* * * * *", fixedFrom, to, 2*time.Minute, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, w := range windows {
		if w.GapMedian <= 0 {
			t.Errorf("expected positive GapMedian, got %v", w.GapMedian)
		}
	}
}

func TestFind_MinFiringsEnforced(t *testing.T) {
	to := fixedFrom.Add(5 * time.Minute)
	// require 100 firings minimum — should return nothing
	windows, err := burst.Find("* * * * *", fixedFrom, to, 2*time.Minute, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) != 0 {
		t.Errorf("expected no windows with high minFirings, got %d", len(windows))
	}
}
