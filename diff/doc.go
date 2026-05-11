// Package diff compares two cron expressions field by field and produces
// human-readable or Markdown-formatted reports of their differences.
//
// # Usage
//
//	r, err := diff.Compare("0 9 * * 1", "0 17 * * 5")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(diff.FormatText(r))
//
// # Output
//
// The result includes a per-field breakdown indicating which fields differ
// between the two expressions, along with a plain-English summary.
//
// Both FormatText and FormatMarkdown are available for rendering results
// suitable for terminal output or documentation respectively.
package diff
