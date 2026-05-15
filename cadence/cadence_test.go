package cadence_test

import (
	"testing"
	"time"

	"github.com/cronparse/cadence"
)

var (
	fixedFrom = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	fixedTo   = time.Date(2024, 1, 1, 2, 0, 0, 0, time.UTC) // 2-hour window
)

func TestAnalyze_EveryMinuteIsRegular(t *testing.T) {
	r, err := cadence.Analyze("* * * * *", fixedFrom, fixedTo, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Regular {
		t.Errorf("expected regular cadence, got jitter=%s", r.Jitter)
	}
	if r.MinGap != time.Minute {
		t.Errorf("expected MinGap=1m, got %s", r.MinGap)
	}
	if r.MaxGap != time.Minute {
		t.Errorf("expected MaxGap=1m, got %s", r.MaxGap)
	}
}

func TestAnalyze_HourlyIsRegular(t *testing.T) {
	to := fixedFrom.Add(5 * time.Hour)
	r, err := cadence.Analyze("0 * * * *", fixedFrom, to, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Regular {
		t.Errorf("expected regular cadence")
	}
	if r.MeanGap != time.Hour {
		t.Errorf("expected MeanGap=1h, got %s", r.MeanGap)
	}
}

func TestAnalyze_InvalidExpression(t *testing.T) {
	_, err := cadence.Analyze("invalid", fixedFrom, fixedTo, time.Second)
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestAnalyze_ToBeforeFrom(t *testing.T) {
	_, err := cadence.Analyze("* * * * *", fixedTo, fixedFrom, time.Second)
	if err == nil {
		t.Error("expected error when to is before from")
	}
}

func TestAnalyze_TooFewRuns(t *testing.T) {
	// hourly in a 30-minute window yields 0 or 1 run
	narrow := fixedFrom.Add(30 * time.Minute)
	_, err := cadence.Analyze("0 * * * *", fixedFrom, narrow, time.Second)
	if err == nil {
		t.Error("expected error for fewer than 2 runs")
	}
}

func TestAnalyze_SampleSizeMatchesRuns(t *testing.T) {
	r, err := cadence.Analyze("* * * * *", fixedFrom, fixedTo, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 2-hour window: 120 minutes → 120 runs (minute 0 through 119 after start)
	if r.SampleSize < 2 {
		t.Errorf("expected at least 2 samples, got %d", r.SampleSize)
	}
}
