package drift

import (
	"fmt"
	"strings"
)

// FormatText returns a human-readable plain-text summary of a drift Result.
func FormatText(r Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Expression : %s\n", r.Expression)
	fmt.Fprintf(&sb, "Samples    : %d intervals\n", r.SampleSize)
	fmt.Fprintf(&sb, "Mean gap   : %s\n", r.MeanGap)
	fmt.Fprintf(&sb, "Min gap    : %s\n", r.MinGap)
	fmt.Fprintf(&sb, "Max gap    : %s\n", r.MaxGap)
	fmt.Fprintf(&sb, "Std dev    : %s\n", r.StdDev)
	fmt.Fprintf(&sb, "Variance   : %.2f ms²\n", r.VarianceMs2)
	return sb.String()
}

// FormatMarkdown returns a Markdown table summarising a drift Result.
func FormatMarkdown(r Result) string {
	var sb strings.Builder
	sb.WriteString("## Drift Report\n\n")
	fmt.Fprintf(&sb, "**Expression:** `%s`\n\n", r.Expression)
	sb.WriteString("| Metric | Value |\n")
	sb.WriteString("|--------|-------|\n")
	fmt.Fprintf(&sb, "| Samples | %d intervals |\n", r.SampleSize)
	fmt.Fprintf(&sb, "| Mean gap | %s |\n", r.MeanGap)
	fmt.Fprintf(&sb, "| Min gap | %s |\n", r.MinGap)
	fmt.Fprintf(&sb, "| Max gap | %s |\n", r.MaxGap)
	fmt.Fprintf(&sb, "| Std dev | %s |\n", r.StdDev)
	fmt.Fprintf(&sb, "| Variance | %.2f ms² |\n", r.VarianceMs2)
	return sb.String()
}
