package digest

import (
	"fmt"
	"strings"
)

// FormatText renders a digest Result as plain text.
func FormatText(r Result) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Expression : %s\n", r.Expression))
	if !r.Valid {
		sb.WriteString(fmt.Sprintf("Valid      : false\n"))
		sb.WriteString(fmt.Sprintf("Error      : %s\n", r.Error))
		return sb.String()
	}
	sb.WriteString(fmt.Sprintf("Valid      : true\n"))
	sb.WriteString(fmt.Sprintf("Description: %s\n", r.Description))
	sb.WriteString(fmt.Sprintf("Frequency  : %s\n", r.Frequency))
	if !r.NextRun.IsZero() {
		sb.WriteString(fmt.Sprintf("Next Run   : %s\n", r.NextRun.UTC().Format("2006-01-02 15:04:05 UTC")))
		sb.WriteString(fmt.Sprintf("Next Run In: %s\n", r.NextRunIn))
	}
	return sb.String()
}

// FormatMarkdown renders a digest Result as a Markdown table.
func FormatMarkdown(r Result) string {
	var sb strings.Builder
	sb.WriteString("| Field | Value |\n")
	sb.WriteString("|---|---|\n")
	sb.WriteString(fmt.Sprintf("| Expression | `%s` |\n", r.Expression))
	if !r.Valid {
		sb.WriteString(fmt.Sprintf("| Valid | ❌ |\n"))
		sb.WriteString(fmt.Sprintf("| Error | %s |\n", r.Error))
		return sb.String()
	}
	sb.WriteString("| Valid | ✅ |\n")
	sb.WriteString(fmt.Sprintf("| Description | %s |\n", r.Description))
	sb.WriteString(fmt.Sprintf("| Frequency | %s |\n", r.Frequency))
	if !r.NextRun.IsZero() {
		sb.WriteString(fmt.Sprintf("| Next Run | %s |\n", r.NextRun.UTC().Format("2006-01-02 15:04:05 UTC")))
		sb.WriteString(fmt.Sprintf("| Next Run In | %s |\n", r.NextRunIn))
	}
	return sb.String()
}
