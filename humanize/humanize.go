// Package humanize provides human-readable descriptions of cron expressions.
package humanize

import (
	"fmt"
	"strings"

	"github.com/user/cronparse/parser"
)

// Describe returns a human-readable description of a cron expression.
func Describe(expr string) (string, error) {
	fields, err := parser.Parse(expr)
	if err != nil {
		return "", fmt.Errorf("invalid expression: %w", err)
	}

	minute := describeField(fields.Minute, "minute", minuteNames)
	hour := describeField(fields.Hour, "hour", hourNames)
	dom := describeField(fields.DayOfMonth, "day of month", nil)
	month := describeField(fields.Month, "month", monthNames)
	dow := describeField(fields.DayOfWeek, "day of week", weekdayNames)

	parts := []string{}

	if fields.Minute == "*" && fields.Hour == "*" {
		parts = append(parts, "every minute")
	} else if fields.Minute == "0" && fields.Hour == "*" {
		parts = append(parts, "at the start of every hour")
	} else {
		parts = append(parts, fmt.Sprintf("at %s past %s", minute, hour))
	}

	if fields.DayOfMonth != "*" {
		parts = append(parts, "on "+dom)
	}
	if fields.Month != "*" {
		parts = append(parts, "in "+month)
	}
	if fields.DayOfWeek != "*" {
		parts = append(parts, "on "+dow)
	}

	return strings.Join(parts, ", "), nil
}

func describeField(val, label string, names map[string]string) string {
	if val == "*" {
		return "every " + label
	}
	if strings.Contains(val, "/") {
		parts := strings.SplitN(val, "/", 2)
		return fmt.Sprintf("every %s %s(s)", parts[1], label)
	}
	if strings.Contains(val, ",") {
		items := strings.Split(val, ",")
		named := make([]string, len(items))
		for i, item := range items {
			named[i] = resolveName(item, names)
		}
		return strings.Join(named, ", ")
	}
	if strings.Contains(val, "-") {
		parts := strings.SplitN(val, "-", 2)
		return fmt.Sprintf("%s through %s", resolveName(parts[0], names), resolveName(parts[1], names))
	}
	return resolveName(val, names)
}

func resolveName(val string, names map[string]string) string {
	if names != nil {
		if name, ok := names[val]; ok {
			return name
		}
	}
	return val
}

var monthNames = map[string]string{
	"1": "January", "2": "February", "3": "March", "4": "April",
	"5": "May", "6": "June", "7": "July", "8": "August",
	"9": "September", "10": "October", "11": "November", "12": "December",
}

var weekdayNames = map[string]string{
	"0": "Sunday", "1": "Monday", "2": "Tuesday", "3": "Wednesday",
	"4": "Thursday", "5": "Friday", "6": "Saturday",
}

var minuteNames map[string]string
var hourNames map[string]string
