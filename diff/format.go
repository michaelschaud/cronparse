package diff

import (
	"fmt"
	"strings"
)

// FormatText renders a Result as a plain-text table.
func FormatText(r *Result) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Diff: %q vs %q\n", r.ExprA, r.ExprB))
	sb.WriteString(strings.Repeat("-", 50) + "\n")
	sb.WriteString(fmt.Sprintf("%-10s %-15s %-15s %s\n", "Field", "A", "B", "Changed"))
	sb.WriteString(strings.Repeat("-", 50) + "\n")
	for _, f := range r.Fields {
		changed := ""
		if !f.Same {
			changed = "*"
		}
		sb.WriteString(fmt.Sprintf("%-10s %-15s %-15s %s\n", f.Field, f.A, f.B, changed))
	}
	sb.WriteString(strings.Repeat("-", 50) + "\n")
	sb.WriteString(Summary(r) + "\n")
	return sb.String()
}

// FormatMarkdown renders a Result as a Markdown table.
func FormatMarkdown(r *Result) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## Cron Diff: `%s` vs `%s`\n\n", r.ExprA, r.ExprB))
	sb.WriteString("| Field | A | B | Changed |\n")
	sb.WriteString("|-------|---|---|---------|\n")
	for _, f := range r.Fields {
		changed := "no"
		if !f.Same {
			changed = "**yes**"
		}
		sb.WriteString(fmt.Sprintf("| %s | `%s` | `%s` | %s |\n", f.Field, f.A, f.B, changed))
	}
	sb.WriteString("\n" + Summary(r) + "\n")
	return sb.String()
}
