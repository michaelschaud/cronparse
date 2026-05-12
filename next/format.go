package next

import (
	"fmt"
	"strings"
	"time"
)

const timeLayout = "2006-01-02 15:04:05 MST"

// FormatText renders a slice of Results as a plain-text report.
func FormatText(results []Result) string {
	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("Expression: %s\n", r.Expression))
		if r.Err != nil {
			sb.WriteString(fmt.Sprintf("  Error: %s\n", r.Err))
		} else {
			for i, t := range r.Runs {
				sb.WriteString(fmt.Sprintf("  [%d] %s\n", i+1, t.Format(timeLayout)))
			}
		}
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// FormatMarkdown renders a slice of Results as a Markdown table per expression.
func FormatMarkdown(results []Result) string {
	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("### `%s`\n\n", r.Expression))
		if r.Err != nil {
			sb.WriteString(fmt.Sprintf("**Error:** %s\n\n", r.Err))
			continue
		}
		sb.WriteString("| # | Next Run |\n")
		sb.WriteString("|---|------------|\n")
		for i, t := range r.Runs {
			sb.WriteString(fmt.Sprintf("| %d | %s |\n", i+1, t.Format(timeLayout)))
		}
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// FormatMergedText renders a merged, deduplicated list of times as plain text.
func FormatMergedText(times []time.Time) string {
	var sb strings.Builder
	sb.WriteString("Merged next runs:\n")
	for i, t := range times {
		sb.WriteString(fmt.Sprintf("  [%d] %s\n", i+1, t.Format(timeLayout)))
	}
	return strings.TrimRight(sb.String(), "\n")
}
