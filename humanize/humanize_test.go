package humanize_test

import (
	"testing"

	"github.com/user/cronparse/humanize"
)

func TestDescribe_EveryMinute(t *testing.T) {
	out, err := humanize.Describe("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "every minute" {
		t.Errorf("got %q, want %q", out, "every minute")
	}
}

func TestDescribe_HourlyAtMinute0(t *testing.T) {
	out, err := humanize.Describe("0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "at the start of every hour" {
		t.Errorf("got %q, want %q", out, "at the start of every hour")
	}
}

func TestDescribe_SpecificTime(t *testing.T) {
	out, err := humanize.Describe("30 9 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "at 30 past 9"
	if out != expected {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestDescribe_WithMonth(t *testing.T) {
	out, err := humanize.Describe("0 0 1 6 *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "at the start of every hour, on 1, in June"
	if out != expected {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestDescribe_WithWeekday(t *testing.T) {
	out, err := humanize.Describe("0 9 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "at the start of every hour, on Monday"
	if out != expected {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestDescribe_StepExpression(t *testing.T) {
	out, err := humanize.Describe("*/15 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "at every 15 minute(s) past every hour"
	if out != expected {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestDescribe_InvalidExpression(t *testing.T) {
	_, err := humanize.Describe("not a cron")
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestDescribe_RangeWeekday(t *testing.T) {
	out, err := humanize.Describe("0 8 * * 1-5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "at the start of every hour, on Monday through Friday"
	if out != expected {
		t.Errorf("got %q, want %q", out, expected)
	}
}
