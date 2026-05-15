package drift_test

import (
	"testing"
	"time"

	"github.com/example/cronparse/drift"
)

func fixedWindow() (time.Time, time.Time) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(2 * time.Hour)
	return from, to
}

func TestMeasure_EveryMinute(t *testing.T) {
	from, to := fixedWindow()
	r, err := drift.Measure("* * * * *", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.SampleSize < 1 {
		t.Error("expected at least one interval")
	}
	// Every-minute schedule has zero drift.
	if r.StdDev != 0 {
		t.Errorf("expected zero std dev for every-minute, got %s", r.StdDev)
	}
	if r.MeanGap != time.Minute {
		t.Errorf("expected mean gap of 1m, got %s", r.MeanGap)
	}
}

func TestMeasure_HourlyZeroDrift(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	r, err := drift.Measure("0 * * * *", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.StdDev != 0 {
		t.Errorf("expected zero std dev for hourly, got %s", r.StdDev)
	}
	if r.MeanGap != time.Hour {
		t.Errorf("expected mean gap of 1h, got %s", r.MeanGap)
	}
}

func TestMeasure_InvalidExpression(t *testing.T) {
	from, to := fixedWindow()
	_, err := drift.Measure("not-a-cron", from, to)
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestMeasure_ToBeforeFrom(t *testing.T) {
	from, to := fixedWindow()
	_, err := drift.Measure("* * * * *", to, from)
	if err == nil {
		t.Error("expected error when to is before from")
	}
}

func TestMeasure_TooFewRuns(t *testing.T) {
	// Yearly expression — unlikely to fire twice in a 2-hour window.
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(2 * time.Hour)
	_, err := drift.Measure("0 0 1 1 *", from, to)
	if err == nil {
		t.Error("expected error for too few runs")
	}
}

func TestMeasure_ResultFieldsPopulated(t *testing.T) {
	from, to := fixedWindow()
	r, err := drift.Measure("* * * * *", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Expression != "* * * * *" {
		t.Errorf("expression not preserved: %q", r.Expression)
	}
	if r.MinGap <= 0 {
		t.Error("min gap should be positive")
	}
	if r.MaxGap < r.MinGap {
		t.Error("max gap should be >= min gap")
	}
}
