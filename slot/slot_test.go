package slot_test

import (
	"testing"
	"time"

	"github.com/cronparse/slot"
)

func fixedFrom() time.Time {
	return time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
}

func TestFind_NoExpressionsAllGap(t *testing.T) {
	from := fixedFrom()
	to := from.Add(2 * time.Hour)

	gaps, err := slot.Find([]string{}, from, to, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gaps) != 1 {
		t.Fatalf("expected 1 gap, got %d", len(gaps))
	}
	if gaps[0].Duration() != 2*time.Hour {
		t.Errorf("expected 2h gap, got %v", gaps[0].Duration())
	}
}

func TestFind_EveryMinuteNoGaps(t *testing.T) {
	from := fixedFrom()
	to := from.Add(30 * time.Minute)

	gaps, err := slot.Find([]string{"* * * * *"}, from, to, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gaps) != 0 {
		t.Errorf("expected no gaps for every-minute schedule, got %d", len(gaps))
	}
}

func TestFind_HourlyHasGaps(t *testing.T) {
	from := fixedFrom() // 09:00
	to := from.Add(2 * time.Hour) // 11:00

	// Fires at :00 of each hour — so 09:00 and 10:00.
	gaps, err := slot.Find([]string{"0 * * * *"}, from, to, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gaps) == 0 {
		t.Fatal("expected gaps for hourly schedule, got none")
	}
	// Each gap should be less than 1 hour.
	for _, g := range gaps {
		if g.Duration() >= time.Hour {
			t.Errorf("gap duration %v unexpectedly large", g.Duration())
		}
	}
}

func TestFind_InvalidExpression(t *testing.T) {
	from := fixedFrom()
	to := from.Add(time.Hour)

	_, err := slot.Find([]string{"not-a-cron"}, from, to, time.Minute)
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}

func TestFind_ToBeforeFrom(t *testing.T) {
	from := fixedFrom()
	to := from.Add(-time.Hour)

	_, err := slot.Find([]string{"* * * * *"}, from, to, time.Minute)
	if err == nil {
		t.Fatal("expected error when to < from")
	}
}

func TestFind_GapDurationIsPositive(t *testing.T) {
	from := fixedFrom()
	to := from.Add(3 * time.Hour)

	gaps, err := slot.Find([]string{"0 * * * *"}, from, to, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, g := range gaps {
		if g.Duration() <= 0 {
			t.Errorf("gap has non-positive duration: %v", g.Duration())
		}
		if !g.End.After(g.Start) {
			t.Errorf("gap end %v not after start %v", g.End, g.Start)
		}
	}
}
