// Package validate provides cron expression validation with detailed error reporting.
package validate

import (
	"fmt"
	"strconv"
	"strings"
)

// FieldSpec defines the allowed range for a cron field.
type FieldSpec struct {
	Name string
	Min  int
	Max  int
}

// fieldSpecs defines the valid ranges for each cron field position.
var fieldSpecs = []FieldSpec{
	{Name: "minute", Min: 0, Max: 59},
	{Name: "hour", Min: 0, Max: 23},
	{Name: "day-of-month", Min: 1, Max: 31},
	{Name: "month", Min: 1, Max: 12},
	{Name: "day-of-week", Min: 0, Max: 6},
}

// ValidationError describes a problem with a specific cron field.
type ValidationError struct {
	Field   string
	Value   string
	Reason  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s field %q: %s", e.Field, e.Value, e.Reason)
}

// Result holds the outcome of validating a cron expression.
type Result struct {
	Valid  bool
	Errors []*ValidationError
}

// Check validates a cron expression and returns a detailed Result.
func Check(expr string) Result {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return Result{
			Valid: false,
			Errors: []*ValidationError{
				{Field: "expression", Value: expr, Reason: fmt.Sprintf("expected 5 fields, got %d", len(fields))},
			},
		}
	}

	var errs []*ValidationError
	for i, field := range fields {
		spec := fieldSpecs[i]
		if err := checkField(field, spec); err != nil {
			errs = append(errs, err)
		}
	}
	return Result{Valid: len(errs) == 0, Errors: errs}
}

func checkField(value string, spec FieldSpec) *ValidationError {
	if value == "*" {
		return nil
	}
	// Handle step values like */5 or 1-5/2
	parts := strings.SplitN(value, "/", 2)
	if len(parts) == 2 {
		step, err := strconv.Atoi(parts[1])
		if err != nil || step < 1 {
			return &ValidationError{Field: spec.Name, Value: value, Reason: "step must be a positive integer"}
		}
	}
	base := parts[0]
	if base == "*" {
		return nil
	}
	// Handle ranges like 1-5
	if strings.Contains(base, "-") {
		rangeParts := strings.SplitN(base, "-", 2)
		lo, err1 := strconv.Atoi(rangeParts[0])
		hi, err2 := strconv.Atoi(rangeParts[1])
		if err1 != nil || err2 != nil {
			return &ValidationError{Field: spec.Name, Value: value, Reason: "range bounds must be integers"}
		}
		if lo > hi {
			return &ValidationError{Field: spec.Name, Value: value, Reason: fmt.Sprintf("range start %d exceeds end %d", lo, hi)}
		}
		if lo < spec.Min || hi > spec.Max {
			return &ValidationError{Field: spec.Name, Value: value, Reason: fmt.Sprintf("range must be within %d-%d", spec.Min, spec.Max)}
		}
		return nil
	}
	// Handle lists like 1,2,3
	for _, item := range strings.Split(base, ",") {
		n, err := strconv.Atoi(item)
		if err != nil {
			return &ValidationError{Field: spec.Name, Value: value, Reason: fmt.Sprintf("%q is not a valid integer", item)}
		}
		if n < spec.Min || n > spec.Max {
			return &ValidationError{Field: spec.Name, Value: value, Reason: fmt.Sprintf("%d out of allowed range %d-%d", n, spec.Min, spec.Max)}
		}
	}
	return nil
}
