package chain_test

import (
	"testing"
	"time"

	"github.com/yourorg/cronparse/chain"
)

var fixedFrom = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestBuild_ValidExpressions(t *testing.T) {
	expressions := map[string]string{
		"every-minute": "* * * * *",
		"hourly":       "0 * * * *",
	}

	result := chain.Build(fixedFrom, 3, expressions)

	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result.Entries))
	}
	for _, e := range result.Entries {
		if e.Err != nil {
			t.Errorf("unexpected error for %q: %v", e.Label, e.Err)
		}
		if len(e.Runs) != 3 {
			t.Errorf("expected 3 runs for %q, got %d", e.Label, len(e.Runs))
		}
	}
}

func TestBuild_MergedIsSorted(t *testing.T) {
	expressions := map[string]string{
		"a": "* * * * *",
		"b": "* * * * *",
	}

	result := chain.Build(fixedFrom, 5, expressions)

	for i := 1; i < len(result.Merged); i++ {
		if result.Merged[i].At.Before(result.Merged[i-1].At) {
			t.Errorf("merged runs not sorted at index %d", i)
		}
	}
}

func TestBuild_CoincidentLabels(t *testing.T) {
	expressions := map[string]string{
		"job-a": "0 * * * *",
		"job-b": "0 * * * *",
	}

	result := chain.Build(fixedFrom, 2, expressions)

	for _, mr := range result.Merged {
		if len(mr.Labels) != 2 {
			t.Errorf("expected 2 labels at %v, got %v", mr.At, mr.Labels)
		}
	}
}

func TestBuild_InvalidExpression(t *testing.T) {
	expressions := map[string]string{
		"bad": "not-a-cron",
		"good": "* * * * *",
	}

	result := chain.Build(fixedFrom, 3, expressions)

	var badEntry *chain.Entry
	for i := range result.Entries {
		if result.Entries[i].Label == "bad" {
			badEntry = &result.Entries[i]
		}
	}
	if badEntry == nil {
		t.Fatal("expected entry for 'bad' label")
	}
	if badEntry.Err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestBuild_EmptyExpressions(t *testing.T) {
	result := chain.Build(fixedFrom, 5, map[string]string{})

	if len(result.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result.Entries))
	}
	if len(result.Merged) != 0 {
		t.Errorf("expected 0 merged runs, got %d", len(result.Merged))
	}
}

func TestBuild_RunsAfterFrom(t *testing.T) {
	expressions := map[string]string{
		"minutely": "* * * * *",
	}

	result := chain.Build(fixedFrom, 4, expressions)

	for _, e := range result.Entries {
		for _, r := range e.Runs {
			if !r.After(fixedFrom) {
				t.Errorf("run %v is not after from time %v", r, fixedFrom)
			}
		}
	}
}
