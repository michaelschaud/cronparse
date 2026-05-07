// Package explain provides a field-by-field breakdown of cron expressions,
// annotating each field (minute, hour, day-of-month, month, day-of-week) with
// a human-readable meaning derived from the parsed values.
//
// It builds on top of the parser package and offers two output formatters:
// plain text (FormatText) and Markdown (FormatMarkdown), suitable for CLI
// display or documentation generation.
//
// Example usage:
//
//	b, err := explain.Explain("*/15 9-17 * * 1-5")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Print(explain.FormatText(b))
package explain
