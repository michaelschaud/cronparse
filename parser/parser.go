package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// CronExpression represents a parsed cron expression with its five fields.
type CronExpression struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
	Raw        string
}

// Parse parses a standard 5-field cron expression string.
// Returns a CronExpression or an error if the expression is invalid.
func Parse(expr string) (*CronExpression, error) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return nil, fmt.Errorf("invalid cron expression: expected 5 fields, got %d", len(fields))
	}

	c := &CronExpression{
		Minute:     fields[0],
		Hour:       fields[1],
		DayOfMonth: fields[2],
		Month:      fields[3],
		DayOfWeek:  fields[4],
		Raw:        expr,
	}

	if err := validateField(c.Minute, 0, 59, "minute"); err != nil {
		return nil, err
	}
	if err := validateField(c.Hour, 0, 23, "hour"); err != nil {
		return nil, err
	}
	if err := validateField(c.DayOfMonth, 1, 31, "day-of-month"); err != nil {
		return nil, err
	}
	if err := validateField(c.Month, 1, 12, "month"); err != nil {
		return nil, err
	}
	if err := validateField(c.DayOfWeek, 0, 6, "day-of-week"); err != nil {
		return nil, err
	}

	return c, nil
}

// validateField checks that a cron field value is syntactically valid
// and within the given min/max bounds.
func validateField(field string, min, max int, name string) error {
	if field == "*" {
		return nil
	}

	// Handle step values like */5 or 1-5/2
	parts := strings.SplitN(field, "/", 2)
	base := parts[0]
	if len(parts) == 2 {
		step, err := strconv.Atoi(parts[1])
		if err != nil || step < 1 {
			return fmt.Errorf("invalid step in %s field: %q", name, field)
		}
	}

	if base == "*" {
		return nil
	}

	// Handle ranges like 1-5
	rangeParts := strings.SplitN(base, "-", 2)
	for _, p := range rangeParts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return fmt.Errorf("invalid value in %s field: %q", name, field)
		}
		if v < min || v > max {
			return fmt.Errorf("%s field value %d out of range [%d-%d]", name, v, min, max)
		}
	}

	return nil
}
