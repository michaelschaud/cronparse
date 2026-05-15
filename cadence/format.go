package cadence

import (
	"fmt"
	"strings"
)

// FormatText returns a plain-text summary of a cadence Result.
func FormatText(r Result) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Expression : %s\n", r.Expression))
	sb.WriteString(fmt.Sprintf("Regular    : %v\n", r.Regular))
	sb.WriteString(fmt.Sprintf("Sample size: %d runs\n", r.SampleSize))
	sb.WriteString(fmt.Sprintf("Min gap    : %s\n", r.MinGap))
	sb.WriteString(fmt.Sprintf("Max gap    : %s\n", r.MaxGap))
	sb.WriteString(fmt.Sprintf("Mean gap   : %s\n", r.MeanGap))
	sb.WriteString(fmt.Sprintf("Jitter     : %s\n", r.Jitter))
	return sb.String()
}

// FormatMarkdown returns a Markdown table summary of a cadence Result.
func FormatMarkdown(r Result) string {
	var sb strings.Builder
	sb.WriteString("| Field | Value |\n")
	sb.WriteString("|---|---|\n")
	sb.WriteString(fmt.Sprintf("| Expression | `%s` |\n", r.Expression))
	sb.WriteString(fmt.Sprintf("| Regular | %v |\n", r.Regular))
	sb.WriteString(fmt.Sprintf("| Sample size | %d runs |\n", r.SampleSize))
	sb.WriteString(fmt.Sprintf("| Min gap | %s |\n", r.MinGap))
	sb.WriteString(fmt.Sprintf("| Max gap | %s |\n", r.MaxGap))
	sb.WriteString(fmt.Sprintf("| Mean gap | %s |\n", r.MeanGap))
	sb.WriteString(fmt.Sprintf("| Jitter | %s |\n", r.Jitter))
	return sb.String()
}
