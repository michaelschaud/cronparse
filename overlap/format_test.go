package overlap_test

import (
	"strings"
	"testing"

	"github.com/yourorg/cronparse/overlap"
)

func TestFormatText_NoOverlaps(t *testing.T) {
	r := &overlap.Result{
		ExpressionA: "0 * * * *",
		ExpressionB: "30 * * * *",
		Overlaps:    nil,
		Count:       0,
	}
	out := overlap.FormatText(r)
	if !strings.Contains(out, "No overlapping") {
		t.Errorf("expected no-overlap message, got: %s", out)
	}
	if !strings.Contains(out, "0 * * * *") {
		t.Errorf("expected expression A in output")
	}
}

func TestFormatText_WithOverlaps(t *testing.T) {
	r, err := overlap.Find("* * * * *", "* * * * *", base, window(3))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := overlap.FormatText(r)
	if !strings.Contains(out, "Times") {
		t.Errorf("expected times section, got: %s", out)
	}
	if r.Count == 0 {
		t.Error("expected overlaps to exist for formatting test")
	}
}

func TestFormatMarkdown_ContainsTable(t *testing.T) {
	r := &overlap.Result{
		ExpressionA: "* * * * *",
		ExpressionB: "* * * * *",
		Overlaps:    nil,
		Count:       0,
	}
	out := overlap.FormatMarkdown(r)
	if !strings.Contains(out, "|") {
		t.Errorf("expected markdown table, got: %s", out)
	}
	if !strings.Contains(out, "## Overlap Report") {
		t.Errorf("expected markdown heading")
	}
}

func TestFormatMarkdown_WithOverlaps(t *testing.T) {
	r, err := overlap.Find("* * * * *", "* * * * *", base, window(3))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := overlap.FormatMarkdown(r)
	if !strings.Contains(out, "Overlapping Times") {
		t.Errorf("expected overlapping times section, got: %s", out)
	}
}
