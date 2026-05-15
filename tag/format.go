package tag

import (
	"fmt"
	"strings"
)

// FormatText renders a Catalog as plain text.
func FormatText(c *Catalog) string {
	var sb strings.Builder
	groups := c.Groups()
	if len(groups) == 0 {
		sb.WriteString("(no entries)\n")
		return sb.String()
	}
	for _, g := range groups {
		fmt.Fprintf(&sb, "[%s]\n", g.Name)
		for _, e := range g.Entries {
			fmt.Fprintf(&sb, "  %-30s %s\n", e.Label, e.Expression)
		}
	}
	return sb.String()
}

// FormatMarkdown renders a Catalog as a Markdown document.
func FormatMarkdown(c *Catalog) string {
	var sb strings.Builder
	groups := c.Groups()
	if len(groups) == 0 {
		sb.WriteString("_(no entries)_\n")
		return sb.String()
	}
	for _, g := range groups {
		fmt.Fprintf(&sb, "## %s\n\n", g.Name)
		sb.WriteString("| Label | Expression |\n")
		sb.WriteString("|-------|------------|\n")
		for _, e := range g.Entries {
			fmt.Fprintf(&sb, "| %s | `%s` |\n", e.Label, e.Expression)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
