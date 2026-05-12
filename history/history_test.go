package history_test

import (
	"testing"
	"time"

	"github.com/cronparse/history"
)

func fixedWindow() (time.Time, time.Time) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	return from, to
}

func TestBetween_EveryMinute(t *testing.T) {
	from, to := fixedWindow()
	r, err := history.Between("* * * * *", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Occurrences) != 60 {
		t.Errorf("expected 60 occurrences, got %d", len(r.Occurrences))
	}
}

func TestBetween_HourlyAtMinute0(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 6, 0, 0, 0, time.UTC)
	r, err := history.Between("0 * * * *", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Occurrences) != 6 {
		t.Errorf("expected 6 occurrences, got %d", len(r.Occurrences))
	}
}

func TestCount_EveryMinute(t *testing.T) {
	from, to := fixedWindow()
	n, err := history.Count("* * * * *", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 60 {
		t.Errorf("expected 60, got %d", n)
	}
}

func TestBetween_InvalidExpression(t *testing.T) {
	from, to := fixedWindow()
	_, err := history.Between("invalid", from, to)
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestBetween_ToBeforeFrom(t *testing.T) {
	from := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := history.Between("* * * * *", from, to)
	if err == nil {
		t.Error("expected error when to is before from")
	}
}

func TestBetween_ResultMetadata(t *testing.T) {
	from, to := fixedWindow()
	r, err := history.Between("0 * * * *", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Expression != "0 * * * *" {
		t.Errorf("expected expression to be preserved, got %q", r.Expression)
	}
	if !r.From.Equal(from) {
		t.Errorf("expected From=%v, got %v", from, r.From)
	}
	if !r.To.Equal(to) {
		t.Errorf("expected To=%v, got %v", to, r.To)
	}
}
