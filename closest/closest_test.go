package closest_test

import (
	"testing"
	"time"

	"github.com/your-org/cronparse/closest"
)

// ref is a fixed Monday 2024-01-15 12:00:00 UTC.
var ref = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestNearest_ReturnsClosest(t *testing.T) {
	// "* * * * *" fires every minute; "0 * * * *" fires at the next full hour.
	exprs := []string{"0 * * * *", "* * * * *"}
	r, err := closest.Nearest(exprs, ref)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Expression != "* * * * *" {
		t.Errorf("expected every-minute expression to be nearest, got %q", r.Expression)
	}
}

func TestFarthest_ReturnsFarthest(t *testing.T) {
	exprs := []string{"0 * * * *", "* * * * *"}
	r, err := closest.Farthest(exprs, ref)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Expression != "0 * * * *" {
		t.Errorf("expected hourly expression to be farthest, got %q", r.Expression)
	}
}

func TestAll_SortedAscending(t *testing.T) {
	exprs := []string{"0 * * * *", "* * * * *", "0 0 * * *"}
	results, err := closest.All(exprs, ref)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for i := 1; i < len(results); i++ {
		if results[i].Next.Before(results[i-1].Next) {
			t.Errorf("results not sorted: index %d (%v) before index %d (%v)",
				i, results[i].Next, i-1, results[i-1].Next)
		}
	}
}

func TestNearest_EmptySlice(t *testing.T) {
	_, err := closest.Nearest([]string{}, ref)
	if err == nil {
		t.Fatal("expected error for empty slice, got nil")
	}
}

func TestFarthest_EmptySlice(t *testing.T) {
	_, err := closest.Farthest(nil, ref)
	if err == nil {
		t.Fatal("expected error for nil slice, got nil")
	}
}

func TestAll_InvalidExpression(t *testing.T) {
	_, err := closest.All([]string{"* * * *"}, ref)
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}
