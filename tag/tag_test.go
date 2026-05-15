package tag_test

import (
	"testing"

	"github.com/yourorg/cronparse/tag"
)

func TestAdd_ValidExpression(t *testing.T) {
	c := tag.NewCatalog()
	if err := c.Add("ops", "hourly", "0 * * * *"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g, ok := c.Group("ops")
	if !ok {
		t.Fatal("group 'ops' not found")
	}
	if len(g.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(g.Entries))
	}
	if g.Entries[0].Label != "hourly" {
		t.Errorf("expected label 'hourly', got %q", g.Entries[0].Label)
	}
}

func TestAdd_InvalidExpression(t *testing.T) {
	c := tag.NewCatalog()
	if err := c.Add("ops", "bad", "not a cron"); err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestGroups_SortedByName(t *testing.T) {
	c := tag.NewCatalog()
	_ = c.Add("zebra", "z", "* * * * *")
	_ = c.Add("alpha", "a", "* * * * *")
	_ = c.Add("mango", "m", "* * * * *")

	groups := c.Groups()
	names := make([]string, len(groups))
	for i, g := range groups {
		names[i] = g.Name
	}
	expected := []string{"alpha", "mango", "zebra"}
	for i, want := range expected {
		if names[i] != want {
			t.Errorf("position %d: want %q, got %q", i, want, names[i])
		}
	}
}

func TestFind_ByLabelSubstring(t *testing.T) {
	c := tag.NewCatalog()
	_ = c.Add("g1", "nightly build", "0 2 * * *")
	_ = c.Add("g1", "weekly release", "0 10 * * 1")
	_ = c.Add("g2", "nightly report", "30 3 * * *")

	results := c.Find("nightly")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestFind_EmptySubstringReturnsAll(t *testing.T) {
	c := tag.NewCatalog()
	_ = c.Add("g", "a", "* * * * *")
	_ = c.Add("g", "b", "0 * * * *")

	results := c.Find("")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestGroup_NotFound(t *testing.T) {
	c := tag.NewCatalog()
	_, ok := c.Group("missing")
	if ok {
		t.Fatal("expected group not to exist")
	}
}
