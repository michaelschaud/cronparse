package closest

import (
	"fmt"
	"strings"
	"time"
)

const timeLayout = "2006-01-02 15:04:05 MST"

// FormatText renders a []Result as a plain-text table.
func FormatText(results []Result) string {
	if len(results) == 0 {
		return "No results."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s  %s\n", "Expression", "Next Run"))
	sb.WriteString(strings.Repeat("-", 55) + "\n")
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("%-30s  %s\n", r.Expression, r.Next.UTC().Format(timeLayout)))
	}
	return sb.String()
}

// FormatMarkdown renders a []Result as a Markdown table.
func FormatMarkdown(results []Result) string {
	if len(results) == 0 {
		return "_No results._"
	}
	var sb strings.Builder
	sb.WriteString("| Expression | Next Run |\n")
	sb.WriteString("|---|---|\n")
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("| `%s` | %s |\n", r.Expression, r.Next.UTC().Format(time.RFC3339)))
	}
	return sb.String()
}
