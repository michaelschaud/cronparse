package heatmap

import (
	"fmt"
	"strings"
	"time"
)

var dayNames = [7]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

// FormatText returns a compact ASCII representation of the heatmap.
// Each row is a day of the week; each column is an hour (0–23).
// Empty cells are shown as ".", non-zero cells show their count (capped at 9).
func FormatText(m *Map) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Heatmap for: %s\n", m.Expression)
	fmt.Fprintf(&sb, "Window: %s – %s\n\n",
		m.From.Format(time.RFC3339), m.To.Format(time.RFC3339))

	// Header
	sb.WriteString("     ")
	for h := 0; h < 24; h++ {
		fmt.Fprintf(&sb, "%02d ", h)
	}
	sb.WriteString("\n")

	for dow := 0; dow < 7; dow++ {
		fmt.Fprintf(&sb, "%s  ", dayNames[dow])
		for h := 0; h < 24; h++ {
			c := m.Grid[dow][h]
			if c == 0 {
				sb.WriteString(".  ")
			} else if c > 9 {
				sb.WriteString("9+ ")
			} else {
				fmt.Fprintf(&sb, "%d  ", c)
			}
		}
		sb.WriteString("\n")
	}

	p := m.Peak()
	if p.Count > 0 {
		fmt.Fprintf(&sb, "\nPeak: %s %02d:xx (%d fires)\n",
			dayNames[p.DayOfWeek], p.Hour, p.Count)
	}
	return sb.String()
}

// FormatMarkdown returns a GitHub-flavoured Markdown table of the heatmap.
func FormatMarkdown(m *Map) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "## Heatmap: `%s`\n\n", m.Expression)
	fmt.Fprintf(&sb, "_Window: %s – %s_\n\n",
		m.From.Format(time.RFC3339), m.To.Format(time.RFC3339))

	// Header row
	sb.WriteString("| Day |")
	for h := 0; h < 24; h++ {
		fmt.Fprintf(&sb, " %02d |", h)
	}
	sb.WriteString("\n|-----|")
	for h := 0; h < 24; h++ {
		sb.WriteString("----|") //nolint:gocritic
	}
	sb.WriteString("\n")

	for dow := 0; dow < 7; dow++ {
		fmt.Fprintf(&sb, "| %s |", dayNames[dow])
		for h := 0; h < 24; h++ {
			c := m.Grid[dow][h]
			if c == 0 {
				sb.WriteString(" . |")
			} else {
				fmt.Fprintf(&sb, " %d |", c)
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
