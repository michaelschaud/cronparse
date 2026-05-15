// Package tag provides a Catalog type for organizing cron expressions
// under human-friendly labels and named groups.
//
// # Overview
//
// A Catalog stores [Entry] values — each pairing a label with a validated
// cron expression — inside named [Group] collections. Entries can be
// retrieved by group name or searched by label substring.
//
// # Example
//
//	c := tag.NewCatalog()
//	_ = c.Add("deployments", "nightly build", "0 2 * * *")
//	_ = c.Add("deployments", "weekly release", "0 10 * * 1")
//	_ = c.Add("monitoring", "health check", "*/5 * * * *")
//
//	fmt.Print(tag.FormatText(c))
package tag
