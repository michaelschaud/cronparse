// Package overlap detects and reports temporal overlaps between two cron
// expressions within a specified time window.
//
// An overlap occurs when both expressions would trigger at the same
// minute boundary. This is useful for identifying scheduling conflicts,
// resource contention, or intentional synchronisation points between jobs.
//
// Basic usage:
//
//	import "github.com/yourorg/cronparse/overlap"
//
//	from := time.Now()
//	to   := from.Add(24 * time.Hour)
//
//	result, err := overlap.Find("0 * * * *", "0 */3 * * *", from, to)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(overlap.FormatText(result))
//
// The window is treated as a half-open interval [from, to).
// Both expressions must be valid standard 5-field cron expressions.
package overlap
