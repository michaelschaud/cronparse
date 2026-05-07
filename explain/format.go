package explain

import (
	"fmt"
	"strings"
)

// FormatText returns a plain-text multi-line representation of a Breakdown.
func FormatText(b *Breakdown) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Expression: %s\n", b.Expression))
	sb.WriteString(strings.Repeat("-", 40) + "\n")
	for _, f := range b.Fields {
		sb.WriteString(fmt.Sprintf("%-16s %-8s %s\n", f.Field+":", f.Value, f.Meaning))
	}
	return sb.String()
}

// FormatMarkdown returns a Markdown table representation of a Breakdown.
func FormatMarkdown(b *Breakdown) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**Expression:** `%s`\n\n", b.Expression))
	sb.WriteString("| Field | Value | Meaning |\n")
	sb.WriteString("|---|---|---|\n")
	for _, f := range b.Fields {
		sb.WriteString(fmt.Sprintf("| %s | `%s` | %s |\n", f.Field, f.Value, f.Meaning))
	}
	return sb.String()
}
