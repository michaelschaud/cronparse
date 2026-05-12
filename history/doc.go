// Package history provides functions for inspecting the historical firing
// pattern of a cron expression over a defined time window.
//
// # Overview
//
// Given a valid cron expression and a [from, to) time range, the history
// package calculates every minute-aligned time at which the expression
// would have matched, returning the results as a slice of [Occurrence]
// values wrapped in a [Result].
//
// # Usage
//
//	import "github.com/cronparse/history"
//
//	from := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
//	to   := time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC)
//
//	r, err := history.Between("0 9 * * 1-5", from, to)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Fired %d times\n", len(r.Occurrences))
//
// Use [Count] when only the total number of firings is needed.
package history
