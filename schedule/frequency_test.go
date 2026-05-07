package schedule_test

import (
	"testing"
	"time"

	"github.com/yourorg/cronparse/schedule"
)

func TestFrequency_EveryMinute(t *testing.T) {
	summary, err := schedule.Frequency("* * * * *", baseTime, 60)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.AvgInterval != time.Minute {
		t.Errorf("expected 1m interval, got %v", summary.AvgInterval)
	}
	if summary.RunsPerHour < 59 || summary.RunsPerHour > 61 {
		t.Errorf("expected ~60 runs/hour, got %.2f", summary.RunsPerHour)
	}
}

func TestFrequency_Hourly(t *testing.T) {
	summary, err := schedule.Frequency("0 * * * *", baseTime, 24)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.AvgInterval != time.Hour {
		t.Errorf("expected 1h interval, got %v", summary.AvgInterval)
	}
	const wantPerDay = 24.0
	if summary.RunsPerDay < wantPerDay-0.1 || summary.RunsPerDay > wantPerDay+0.1 {
		t.Errorf("expected ~24 runs/day, got %.2f", summary.RunsPerDay)
	}
}

func TestFrequency_InvalidExpression(t *testing.T) {
	_, err := schedule.Frequency("not valid", baseTime, 10)
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestFrequency_SmallSampleClamped(t *testing.T) {
	// sampleSize < 2 should be clamped to 2 internally
	summary, err := schedule.Frequency("* * * * *", baseTime, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.AvgInterval <= 0 {
		t.Error("expected positive average interval")
	}
}
