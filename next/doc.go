// Package next provides multi-expression next-run scheduling.
//
// It builds on the forecast package to support querying multiple cron
// expressions simultaneously, returning per-expression results or a
// merged, deduplicated, chronologically sorted list of upcoming run times.
//
// # Basic usage
//
//	results := next.ForExpressions(
//		[]string{"* * * * *", "0 9 * * 1"},
//		time.Now(),
//		5,
//	)
//	for _, r := range results {
//		if r.Err != nil {
//			log.Println(r.Expression, r.Err)
//			continue
//		}
//		for _, t := range r.Runs {
//			fmt.Println(r.Expression, t)
//		}
//	}
//
// # Merged output
//
// Use Merged to combine runs from all expressions into a single sorted,
// deduplicated slice with an optional upper limit:
//
//	times, err := next.Merged(exprs, time.Now(), 10, 20)
//
// # Formatting
//
// FormatText and FormatMarkdown render Results slices for CLI or
// documentation output. FormatMergedText renders a merged time list.
package next
