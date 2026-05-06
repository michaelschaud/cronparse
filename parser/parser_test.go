package parser

import (
	"testing"
)

func TestParse_ValidExpressions(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"every minute", "* * * * *"},
		{"specific time", "30 9 * * 1"},
		{"every 5 minutes", "*/5 * * * *"},
		{"range", "0 9-17 * * 1-5"},
		{"step with range", "0 1-23/2 * * *"},
		{"midnight daily", "0 0 * * *"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Parse(tt.expr)
			if err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.expr, err)
				return
			}
			if c.Raw != tt.expr {
				t.Errorf("expected Raw = %q, got %q", tt.expr, c.Raw)
			}
		})
	}
}

func TestParse_InvalidExpressions(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"too few fields", "* * * *"},
		{"too many fields", "* * * * * *"},
		{"minute out of range", "60 * * * *"},
		{"hour out of range", "* 24 * * *"},
		{"day-of-month out of range", "* * 32 * *"},
		{"month out of range", "* * * 13 *"},
		{"dow out of range", "* * * * 7"},
		{"invalid step", "*/0 * * * *"},
		{"non-numeric", "abc * * * *"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.expr)
			if err == nil {
				t.Errorf("Parse(%q) expected error, got nil", tt.expr)
			}
		})
	}
}

func TestParse_FieldValues(t *testing.T) {
	c, err := Parse("30 9 15 6 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Minute != "30" {
		t.Errorf("expected Minute=30, got %s", c.Minute)
	}
	if c.Hour != "9" {
		t.Errorf("expected Hour=9, got %s", c.Hour)
	}
	if c.DayOfMonth != "15" {
		t.Errorf("expected DayOfMonth=15, got %s", c.DayOfMonth)
	}
	if c.Month != "6" {
		t.Errorf("expected Month=6, got %s", c.Month)
	}
	if c.DayOfWeek != "1" {
		t.Errorf("expected DayOfWeek=1, got %s", c.DayOfWeek)
	}
}
