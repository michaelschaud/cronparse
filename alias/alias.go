// Package alias provides named shortcuts for common cron expressions.
// Users can register human-friendly names (e.g. "daily", "weekly") and
// resolve them back to standard five-field cron strings.
package alias

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yourorg/cronparse/validate"
)

// Builtin aliases that ship with the package.
var builtins = map[string]string{
	"yearly":   "0 0 1 1 *",
	"annually": "0 0 1 1 *",
	"monthly":  "0 0 1 * *",
	"weekly":   "0 0 * * 0",
	"daily":    "0 0 * * *",
	"midnight": "0 0 * * *",
	"hourly":   "0 * * * *",
}

// Registry holds a set of named cron expression aliases.
type Registry struct {
	entries map[string]string
}

// NewRegistry returns a Registry pre-loaded with built-in aliases.
func NewRegistry() *Registry {
	r := &Registry{entries: make(map[string]string)}
	for k, v := range builtins {
		r.entries[k] = v
	}
	return r
}

// Add registers a new alias. The expression is validated before storage.
// Returns an error if the name is empty, already registered, or the
// expression is invalid.
func (r *Registry) Add(name, expr string) error {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return fmt.Errorf("alias name must not be empty")
	}
	if _, exists := r.entries[name]; exists {
		return fmt.Errorf("alias %q is already registered", name)
	}
	if err := validate.Check(expr); err != nil {
		return fmt.Errorf("invalid expression for alias %q: %w", name, err)
	}
	r.entries[name] = expr
	return nil
}

// Resolve returns the cron expression for the given alias name.
// If the name is not found, it is returned as-is so callers can
// transparently accept both aliases and raw expressions.
func (r *Registry) Resolve(name string) string {
	if expr, ok := r.entries[strings.ToLower(strings.TrimSpace(name))]; ok {
		return expr
	}
	return name
}

// List returns all registered aliases sorted alphabetically.
func (r *Registry) List() []Entry {
	out := make([]Entry, 0, len(r.entries))
	for name, expr := range r.entries {
		out = append(out, Entry{Name: name, Expression: expr})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Entry is a single alias name/expression pair.
type Entry struct {
	Name       string
	Expression string
}
