// Package diff provides utilities for comparing two cron expressions
// and describing the differences between their schedules in human-readable form.
package diff

import (
	"fmt"
	"strings"

	"github.com/yourorg/cronparse/parser"
)

// FieldDiff describes the difference between a single cron field.
type FieldDiff struct {
	Field string
	A     string
	B     string
	Same  bool
}

// Result holds the full diff between two cron expressions.
type Result struct {
	ExprA  string
	ExprB  string
	Fields []FieldDiff
	Same   bool
}

var fieldNames = []string{"Minute", "Hour", "Day", "Month", "Weekday"}

// Compare returns a Result describing the field-by-field differences
// between cron expression a and b. Returns an error if either expression
// is invalid.
func Compare(a, b string) (*Result, error) {
	pa, err := parser.Parse(a)
	if err != nil {
		return nil, fmt.Errorf("expression A: %w", err)
	}
	pb, err := parser.Parse(b)
	if err != nil {
		return nil, fmt.Errorf("expression B: %w", err)
	}

	fieldsA := []string{pa.Minute, pa.Hour, pa.Day, pa.Month, pa.Weekday}
	fieldsB := []string{pb.Minute, pb.Hour, pb.Day, pb.Month, pb.Weekday}

	allSame := true
	diffs := make([]FieldDiff, len(fieldNames))
	for i, name := range fieldNames {
		same := fieldsA[i] == fieldsB[i]
		if !same {
			allSame = false
		}
		diffs[i] = FieldDiff{
			Field: name,
			A:     fieldsA[i],
			B:     fieldsB[i],
			Same:  same,
		}
	}

	return &Result{
		ExprA:  a,
		ExprB:  b,
		Fields: diffs,
		Same:   allSame,
	}, nil
}

// Summary returns a short human-readable summary of the diff result.
func Summary(r *Result) string {
	if r.Same {
		return fmt.Sprintf("Expressions %q and %q are identical.", r.ExprA, r.ExprB)
	}
	var changed []string
	for _, f := range r.Fields {
		if !f.Same {
			changed = append(changed, strings.ToLower(f.Field))
		}
	}
	return fmt.Sprintf("Expressions differ in: %s.", strings.Join(changed, ", "))
}
