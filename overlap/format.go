package overlap

import (
	"fmt"
	"strings"
	"time"
)

const timeLayout = "2006-01-02 15:04"

// FormatText returns a plain-text summary of the overlap result.
func FormatText(r *Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Overlap report\n")
	fmt.Fprintf(&sb, "  Expression A : %s\n", r.ExpressionA)
	fmt.Fprintf(&sb, "  Expression B : %s\n", r.ExpressionB)
	fmt.Fprintf(&sb, "  Total overlaps: %d\n", r.Count)
	if r.Count == 0 {
		sb.WriteString("  No overlapping fire times found.\n")
		return sb.String()
	}
	sb.WriteString("  Times (UTC):\n")
	for _, t := range r.Overlaps {
		fmt.Fprintf(&sb, "    - %s\n", t.UTC().Format(timeLayout))
	}
	return sb.String()
}

// FormatMarkdown returns a Markdown-formatted summary of the overlap result.
func FormatMarkdown(r *Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "## Overlap Report\n\n")
	fmt.Fprintf(&sb, "| Field | Value |\n|---|---|\n")
	fmt.Fprintf(&sb, "| Expression A | `%s` |\n", r.ExpressionA)
	fmt.Fprintf(&sb, "| Expression B | `%s` |\n", r.ExpressionB)
	fmt.Fprintf(&sb, "| Total Overlaps | %d |\n\n", r.Count)
	if r.Count == 0 {
		sb.WriteString("_No overlapping fire times found._\n")
		return sb.String()
	}
	sb.WriteString("### Overlapping Times (UTC)\n\n")
	for _, t := range r.Overlaps {
		fmt.Fprintf(&sb, "- `%s`\n", t.UTC().Format(timeLayout))
	}
	return sb.String()
}

// formatTime is a helper used in tests.
func formatTime(t time.Time) string {
	return t.UTC().Format(timeLayout)
}
