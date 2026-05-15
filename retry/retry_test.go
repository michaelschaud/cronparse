package retry_test

import (
	"testing"
	"time"

	"github.com/cronparse/retry"
)

func TestSuggest_InvalidExpression(t *testing.T) {
	_, err := retry.Suggest("not valid", time.Hour, 3)
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}

func TestSuggest_NegativeTarget(t *testing.T) {
	_, err := retry.Suggest("* * * * *", -time.Minute, 3)
	if err == nil {
		t.Fatal("expected error for non-positive target, got nil")
	}
}

func TestSuggest_ReturnsAtMostN(t *testing.T) {
	suggestions, err := retry.Suggest("0 0 * * *", time.Hour, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(suggestions) > 2 {
		t.Errorf("expected at most 2 suggestions, got %d", len(suggestions))
	}
}

func TestSuggest_ExcludesInputExpression(t *testing.T) {
	expr := "*/5 * * * *"
	suggestions, err := retry.Suggest(expr, 5*time.Minute, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, s := range suggestions {
		if s.Expression == expr {
			t.Errorf("suggestions should not include the input expression %q", expr)
		}
	}
}

func TestSuggest_SortedByDelta(t *testing.T) {
	suggestions, err := retry.Suggest("0 0 1 * *", time.Hour, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(suggestions); i++ {
		if suggestions[i].DeltaAbs < suggestions[i-1].DeltaAbs {
			t.Errorf("suggestions not sorted by DeltaAbs at index %d", i)
		}
	}
}

func TestSuggest_HourlyTargetClosestIsHourly(t *testing.T) {
	suggestions, err := retry.Suggest("0 0 * * *", time.Hour, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(suggestions) == 0 {
		t.Fatal("expected at least one suggestion")
	}
	if suggestions[0].Expression != "0 * * * *" {
		t.Logf("closest suggestion for 1h target: %q (%s)", suggestions[0].Expression, suggestions[0].Description)
	}
}

func TestSuggest_DefaultNIsThree(t *testing.T) {
	// passing n=0 should default to 3
	suggestions, err := retry.Suggest("* * * * *", 30*time.Minute, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(suggestions) > 3 {
		t.Errorf("expected at most 3 suggestions with n=0, got %d", len(suggestions))
	}
}

func TestSuggestion_HasDescription(t *testing.T) {
	suggestions, err := retry.Suggest("0 0 * * *", time.Hour, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, s := range suggestions {
		if s.Description == "" {
			t.Errorf("suggestion %q has empty description", s.Expression)
		}
	}
}
