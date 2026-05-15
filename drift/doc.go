// Package drift analyses the temporal consistency of cron expressions by
// computing statistics over the gaps between consecutive scheduled runs within
// a given time window.
//
// A perfectly regular expression (e.g. "* * * * *") will have zero standard
// deviation and zero variance, while expressions that fire at irregular
// intervals (e.g. on specific days of the week combined with specific months)
// will exhibit measurable drift.
//
// Usage:
//
//	from := time.Now()
//	to := from.Add(7 * 24 * time.Hour)
//
//	r, err := drift.Measure("0 9 * * 1-5", from, to)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Print(drift.FormatText(r))
package drift
