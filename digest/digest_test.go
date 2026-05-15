package digest_test

import (
	"testing"
	"time"

	"github.com/cronparse/digest"
)

var fixedFrom = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func TestOf_ValidExpression(t *testing.T) {
	r := digest.Of("* * * * *", fixedFrom)
	if !r.Valid {
		t.Fatalf("expected valid, got error: %s", r.Error)
	}
	if r.Expression != "* * * * *" {
		t.Errorf("unexpected expression: %s", r.Expression)
	}
	if r.Description == "" {
		t.Error("expected non-empty description")
	}
	if r.Frequency == "" {
		t.Error("expected non-empty frequency label")
	}
	if r.NextRun.IsZero() {
		t.Error("expected non-zero next run")
	}
	if r.NextRun.Before(fixedFrom) {
		t.Errorf("next run %v is before from %v", r.NextRun, fixedFrom)
	}
}

func TestOf_InvalidExpression(t *testing.T) {
	r := digest.Of("invalid", fixedFrom)
	if r.Valid {
		t.Fatal("expected invalid")
	}
	if r.Error == "" {
		t.Error("expected non-empty error")
	}
	if r.Description != "" {
		t.Errorf("expected empty description for invalid expr, got %q", r.Description)
	}
}

func TestOf_HourlyFrequencyLabel(t *testing.T) {
	r := digest.Of("0 * * * *", fixedFrom)
	if !r.Valid {
		t.Fatalf("expected valid: %s", r.Error)
	}
	if r.Frequency != "hourly" {
		t.Errorf("expected 'hourly', got %q", r.Frequency)
	}
}

func TestOf_EveryMinuteFrequencyLabel(t *testing.T) {
	r := digest.Of("* * * * *", fixedFrom)
	if !r.Valid {
		t.Fatalf("expected valid: %s", r.Error)
	}
	if r.Frequency != "every minute" {
		t.Errorf("expected 'every minute', got %q", r.Frequency)
	}
}

func TestOf_NextRunInIsPositive(t *testing.T) {
	r := digest.Of("* * * * *", fixedFrom)
	if !r.Valid {
		t.Fatalf("expected valid: %s", r.Error)
	}
	if r.NextRunIn <= 0 {
		t.Errorf("expected positive NextRunIn, got %v", r.NextRunIn)
	}
}
