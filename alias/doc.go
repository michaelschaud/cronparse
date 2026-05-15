// Package alias provides a Registry for mapping human-friendly shorthand
// names to standard five-field cron expressions.
//
// # Built-in aliases
//
// The following aliases are available out of the box:
//
//	"yearly"   → "0 0 1 1 *"
//	"annually" → "0 0 1 1 *"
//	"monthly"  → "0 0 1 * *"
//	"weekly"   → "0 0 * * 0"
//	"daily"    → "0 0 * * *"
//	"midnight" → "0 0 * * *"
//	"hourly"   → "0 * * * *"
//
// # Usage
//
//	r := alias.NewRegistry()
//
//	// Register a custom alias
//	r.Add("every5m", "*/5 * * * *")
//
//	// Resolve transparently accepts aliases and raw expressions
//	expr := r.Resolve("daily")   // → "0 0 * * *"
//	expr  = r.Resolve("*/5 * * * *") // → "*/5 * * * *" (passthrough)
//
//	// Enumerate all registered aliases
//	for _, e := range r.List() {
//		fmt.Printf("%s → %s\n", e.Name, e.Expression)
//	}
package alias
