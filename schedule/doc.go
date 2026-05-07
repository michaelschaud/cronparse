// Package schedule provides higher-level schedule analysis utilities built on
// top of the parser and forecast packages.
//
// It includes:
//
//   - Compare: diff two cron expressions to find overlapping and exclusive run times.
//   - Frequency: estimate how often a cron expression fires (runs per hour/day, average interval).
//
// Example usage:
//
//	base := time.Now()
//	diff, err := schedule.Compare("*/5 * * * *", "*/10 * * * *", base, 20)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Common runs:", diff.Common)
//
//	summary, err := schedule.Frequency("0 * * * *", base, 48)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Runs per day: %.1f\n", summary.RunsPerDay)
package schedule
