package alias_test

import (
	"testing"

	"github.com/yourorg/cronparse/alias"
)

func TestResolve_BuiltinAlias(t *testing.T) {
	r := alias.NewRegistry()
	expr := r.Resolve("daily")
	if expr != "0 0 * * *" {
		t.Errorf("expected '0 0 * * *', got %q", expr)
	}
}

func TestResolve_UnknownReturnsAsIs(t *testing.T) {
	r := alias.NewRegistry()
	raw := "*/5 * * * *"
	if got := r.Resolve(raw); got != raw {
		t.Errorf("expected passthrough %q, got %q", raw, got)
	}
}

func TestResolve_CaseInsensitive(t *testing.T) {
	r := alias.NewRegistry()
	if r.Resolve("HOURLY") != r.Resolve("hourly") {
		t.Error("resolution should be case-insensitive")
	}
}

func TestAdd_ValidAlias(t *testing.T) {
	r := alias.NewRegistry()
	if err := r.Add("every5m", "*/5 * * * *"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := r.Resolve("every5m"); got != "*/5 * * * *" {
		t.Errorf("expected '*/5 * * * *', got %q", got)
	}
}

func TestAdd_InvalidExpression(t *testing.T) {
	r := alias.NewRegistry()
	if err := r.Add("bad", "not a cron"); err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestAdd_EmptyName(t *testing.T) {
	r := alias.NewRegistry()
	if err := r.Add("", "* * * * *"); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestAdd_DuplicateName(t *testing.T) {
	r := alias.NewRegistry()
	if err := r.Add("snap", "* * * * *"); err != nil {
		t.Fatalf("first add failed: %v", err)
	}
	if err := r.Add("snap", "0 * * * *"); err == nil {
		t.Error("expected error on duplicate alias name")
	}
}

func TestList_ContainsBuiltins(t *testing.T) {
	r := alias.NewRegistry()
	entries := r.List()
	names := make(map[string]bool)
	for _, e := range entries {
		names[e.Name] = true
	}
	for _, want := range []string{"daily", "hourly", "weekly", "monthly", "yearly"} {
		if !names[want] {
			t.Errorf("expected builtin alias %q in list", want)
		}
	}
}

func TestList_SortedAlphabetically(t *testing.T) {
	r := alias.NewRegistry()
	entries := r.List()
	for i := 1; i < len(entries); i++ {
		if entries[i].Name < entries[i-1].Name {
			t.Errorf("list not sorted: %q before %q", entries[i-1].Name, entries[i].Name)
		}
	}
}
