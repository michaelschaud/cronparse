package validate_test

import (
	"testing"

	"github.com/yourorg/cronparse/validate"
)

func TestCheck_ValidExpressions(t *testing.T) {
	cases := []string{
		"* * * * *",
		"0 * * * *",
		"0 12 * * *",
		"*/5 * * * *",
		"0 9-17 * * 1-5",
		"30 6 1,15 * *",
		"0 0 * 1 0",
		"*/10 */2 * * *",
	}
	for _, expr := range cases {
		t.Run(expr, func(t *testing.T) {
			res := validate.Check(expr)
			if !res.Valid {
				t.Errorf("expected valid, got errors: %v", res.Errors)
			}
		})
	}
}

func TestCheck_WrongFieldCount(t *testing.T) {
	cases := []string{
		"",
		"* * *",
		"* * * * * *",
	}
	for _, expr := range cases {
		t.Run(expr, func(t *testing.T) {
			res := validate.Check(expr)
			if res.Valid {
				t.Error("expected invalid")
			}
			if len(res.Errors) == 0 {
				t.Error("expected at least one error")
			}
		})
	}
}

func TestCheck_OutOfRangeValues(t *testing.T) {
	cases := []struct {
		expr  string
		field string
	}{
		{"60 * * * *", "minute"},
		{"* 24 * * *", "hour"},
		{"* * 0 * *", "day-of-month"},
		{"* * 32 * *", "day-of-month"},
		{"* * * 0 *", "month"},
		{"* * * 13 *", "month"},
		{"* * * * 7", "day-of-week"},
	}
	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			res := validate.Check(tc.expr)
			if res.Valid {
				t.Errorf("expected invalid for %q", tc.expr)
			}
			found := false
			for _, e := range res.Errors {
				if e.Field == tc.field {
					found = true
				}
			}
			if !found {
				t.Errorf("expected error for field %q, got %v", tc.field, res.Errors)
			}
		})
	}
}

func TestCheck_InvalidStep(t *testing.T) {
	cases := []string{
		"*/0 * * * *",
		"*/abc * * * *",
	}
	for _, expr := range cases {
		t.Run(expr, func(t *testing.T) {
			res := validate.Check(expr)
			if res.Valid {
				t.Errorf("expected invalid for %q", expr)
			}
		})
	}
}

func TestCheck_InvalidRange(t *testing.T) {
	cases := []string{
		"5-3 * * * *",
		"* 20-25 * * *",
	}
	for _, expr := range cases {
		t.Run(expr, func(t *testing.T) {
			res := validate.Check(expr)
			if res.Valid {
				t.Errorf("expected invalid for %q", expr)
			}
		})
	}
}
