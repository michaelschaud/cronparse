package forecast_test

import (
	"testing"
	"time"

	"github.com/cronparse/forecast"
)

func fixedTime(year, month, day, hour, min int) time.Time {
	return time.Date(year, time.Month(month), day, hour, min, 0, 0, time.UTC)
}

func TestNextRun_EveryMinute(t *testing.T) {
	start := fixedTime(2024, 1, 15, 10, 30)
	next, err := forecast.NextRun("* * * * *", start)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := fixedTime(2024, 1, 15, 10, 31)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestNextRun_HourlyAtMinute0(t *testing.T) {
	start := fixedTime(2024, 1, 15, 10, 30)
	next, err := forecast.NextRun("0 * * * *", start)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := fixedTime(2024, 1, 15, 11, 0)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestNextRuns_ReturnsNResults(t *testing.T) {
	start := fixedTime(2024, 1, 15, 10, 0)
	runs, err := forecast.NextRuns("0 * * * *", start, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(runs) != 3 {
		t.Fatalf("expected 3 runs, got %d", len(runs))
	}
	expected := []time.Time{
		fixedTime(2024, 1, 15, 11, 0),
		fixedTime(2024, 1, 15, 12, 0),
		fixedTime(2024, 1, 15, 13, 0),
	}
	for i, r := range runs {
		if !r.Equal(expected[i]) {
			t.Errorf("run[%d]: expected %v, got %v", i, expected[i], r)
		}
	}
}

func TestNextRuns_InvalidExpression(t *testing.T) {
	start := fixedTime(2024, 1, 15, 10, 0)
	_, err := forecast.NextRuns("invalid", start, 1)
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestNextRuns_InvalidN(t *testing.T) {
	start := fixedTime(2024, 1, 15, 10, 0)
	_, err := forecast.NextRuns("* * * * *", start, 0)
	if err == nil {
		t.Error("expected error for n=0, got nil")
	}
}
