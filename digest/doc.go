// Package digest provides a high-level summary of a cron expression,
// combining validation, human-readable description, frequency classification,
// and next-run forecasting into a single Result value.
//
// # Usage
//
//	import "github.com/cronparse/digest"
//
//	r := digest.Of("0 9 * * 1-5", time.Now())
//	if r.Valid {
//	    fmt.Println(r.Description)  // "At 09:00, Monday through Friday"
//	    fmt.Println(r.Frequency)    // "daily"
//	    fmt.Println(r.NextRun)      // next scheduled time
//	}
//
// # Formatting
//
// Use FormatText or FormatMarkdown to render a Result for display:
//
//	fmt.Print(digest.FormatText(r))
//	fmt.Print(digest.FormatMarkdown(r))
package digest
