// Package burst detects time windows where a cron expression fires in rapid
// succession — i.e., where consecutive inter-firing gaps fall below a given
// threshold for at least a minimum number of events.
//
// # Overview
//
// Some cron expressions, particularly those combining step values or ranges,
// can produce uneven firing distributions. The burst package surfaces these
// clusters so callers can reason about load spikes or overlapping jobs.
//
// # Usage
//
//	windows, err := burst.Find(
//		"*/5 * * * *",
//		time.Now(),
//		time.Now().Add(2*time.Hour),
//		10*time.Minute,
//		3,
//	)
//
// Each returned [Window] contains the slice of firing times, the start/end
// of the burst, and the median inter-firing gap within that window.
package burst
