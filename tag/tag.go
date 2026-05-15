// Package tag provides utilities for labeling cron expressions with
// human-friendly tags and grouping them by category.
package tag

import (
	"fmt"
	"sort"

	"github.com/yourorg/cronparse/parser"
)

// Entry associates a cron expression with a user-defined label.
type Entry struct {
	Label      string
	Expression string
}

// Group holds a named collection of tagged cron entries.
type Group struct {
	Name    string
	Entries []Entry
}

// Catalog holds all groups and provides lookup operations.
type Catalog struct {
	groups map[string]*Group
}

// NewCatalog creates an empty Catalog.
func NewCatalog() *Catalog {
	return &Catalog{groups: make(map[string]*Group)}
}

// Add inserts a labeled expression into the named group.
// Returns an error if the expression is invalid.
func (c *Catalog) Add(groupName, label, expr string) error {
	if _, err := parser.Parse(expr); err != nil {
		return fmt.Errorf("tag: invalid expression %q: %w", expr, err)
	}
	g, ok := c.groups[groupName]
	if !ok {
		g = &Group{Name: groupName}
		c.groups[groupName] = g
	}
	g.Entries = append(g.Entries, Entry{Label: label, Expression: expr})
	return nil
}

// Group returns the named group and whether it exists.
func (c *Catalog) Group(name string) (*Group, bool) {
	g, ok := c.groups[name]
	return g, ok
}

// Groups returns all groups sorted by name.
func (c *Catalog) Groups() []*Group {
	out := make([]*Group, 0, len(c.groups))
	for _, g := range c.groups {
		out = append(out, g)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Find returns all entries across all groups whose label contains substr.
func (c *Catalog) Find(substr string) []Entry {
	var results []Entry
	for _, g := range c.Groups() {
		for _, e := range g.Entries {
			if contains(e.Label, substr) {
				results = append(results, e)
			}
		}
	}
	return results
}

func contains(s, sub string) bool {
	return len(sub) == 0 || len(s) >= len(sub) && (s == sub || len(s) > 0 && containsRune(s, sub))
}

func containsRune(s, sub string) bool {
	for i := range s {
		if i+len(sub) <= len(s) && s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
