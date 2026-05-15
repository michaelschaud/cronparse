// Package cadence measures the temporal regularity of a cron expression.
//
// Given a cron expression and a time window, Analyze computes the gaps between
// consecutive scheduled runs and reports whether the schedule fires at a
// consistent interval. A configurable jitter threshold controls what counts as
// "regular".
//
// # Example
//
//	result, err := cadence.Analyze("*/5 * * * *", from, to, time.Second)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(cadence.FormatText(result))
//
// The returned Result includes MinGap, MaxGap, MeanGap, Jitter, and a boolean
// Regular flag indicating whether jitter is within the supplied threshold.
package cadence
