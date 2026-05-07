// Package explain provides field-by-field breakdown of cron expressions.
package explain

import (
	"fmt"
	"strings"

	"github.com/yourorg/cronparse/parser"
)

// FieldExplanation holds the name and human-readable description of a single cron field.
type FieldExplanation struct {
	Field string
	Value string
	Meaning string
}

// Breakdown holds the full explanation of a parsed cron expression.
type Breakdown struct {
	Expression string
	Fields []FieldExplanation
}

var fieldNames = []string{"Minute", "Hour", "Day-of-Month", "Month", "Day-of-Week"}

var monthNames = []string{"", "January", "February", "March", "April", "May",
	"June", "July", "August", "September", "October", "November", "December"}

var weekdayNames = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

// Explain parses a cron expression and returns a field-by-field breakdown.
func Explain(expr string) (*Breakdown, error) {
	cron, err := parser.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("explain: %w", err)
	}

	rawFields := strings.Fields(expr)
	values := [][]int{cron.Minutes, cron.Hours, cron.DaysOfMonth, cron.Months, cron.DaysOfWeek}
	ranges := [][2]int{{0, 59}, {0, 23}, {1, 31}, {1, 12}, {0, 6}}

	fields := make([]FieldExplanation, 5)
	for i := 0; i < 5; i++ {
		fields[i] = FieldExplanation{
			Field:   fieldNames[i],
			Value:   rawFields[i],
			Meaning: describeParsed(i, rawFields[i], values[i], ranges[i]),
		}
	}

	return &Breakdown{Expression: expr, Fields: fields}, nil
}

func describeParsed(fieldIdx int, raw string, vals []int, r [2]int) string {
	if raw == "*" {
		return fmt.Sprintf("every %s", strings.ToLower(fieldNames[fieldIdx]))
	}
	if strings.HasPrefix(raw, "*/") {
		step := raw[2:]
		return fmt.Sprintf("every %s %s(s)", step, strings.ToLower(fieldNames[fieldIdx]))
	}
	named := make([]string, len(vals))
	for i, v := range vals {
		named[i] = resolveFieldName(fieldIdx, v)
	}
	return strings.Join(named, ", ")
}

func resolveFieldName(fieldIdx, val int) string {
	switch fieldIdx {
	case 3:
		if val >= 1 && val <= 12 {
			return monthNames[val]
		}
	case 4:
		if val >= 0 && val <= 6 {
			return weekdayNames[val]
		}
	}
	return fmt.Sprintf("%d", val)
}
