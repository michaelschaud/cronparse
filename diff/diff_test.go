package diff_test

import (
	"strings"
	"testing"

	"github.com/yourorg/cronparse/diff"
)

func TestCompare_SameExpression(t *testing.T) {
	r, err := diff.Compare("0 * * * *", "0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Same {
		t.Error("expected expressions to be the same")
	}
	for _, f := range r.Fields {
		if !f.Same {
			t.Errorf("field %s should be the same", f.Field)
		}
	}
}

func TestCompare_DifferentMinute(t *testing.T) {
	r, err := diff.Compare("0 * * * *", "30 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Same {
		t.Error("expected expressions to differ")
	}
	if r.Fields[0].Same {
		t.Error("expected Minute field to differ")
	}
	for _, f := range r.Fields[1:] {
		if !f.Same {
			t.Errorf("field %s should be the same", f.Field)
		}
	}
}

func TestCompare_MultipleFieldsDiffer(t *testing.T) {
	r, err := diff.Compare("0 9 * * 1", "30 18 * * 5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Same {
		t.Error("expected expressions to differ")
	}
	changedCount := 0
	for _, f := range r.Fields {
		if !f.Same {
			changedCount++
		}
	}
	if changedCount != 3 {
		t.Errorf("expected 3 changed fields, got %d", changedCount)
	}
}

func TestCompare_InvalidExpressionA(t *testing.T) {
	_, err := diff.Compare("invalid", "* * * * *")
	if err == nil {
		t.Error("expected error for invalid expression A")
	}
	if !strings.Contains(err.Error(), "expression A") {
		t.Errorf("error should mention expression A, got: %v", err)
	}
}

func TestCompare_InvalidExpressionB(t *testing.T) {
	_, err := diff.Compare("* * * * *", "bad expr")
	if err == nil {
		t.Error("expected error for invalid expression B")
	}
	if !strings.Contains(err.Error(), "expression B") {
		t.Errorf("error should mention expression B, got: %v", err)
	}
}

func TestSummary_Same(t *testing.T) {
	r, _ := diff.Compare("* * * * *", "* * * * *")
	s := diff.Summary(r)
	if !strings.Contains(s, "identical") {
		t.Errorf("expected 'identical' in summary, got: %s", s)
	}
}

func TestSummary_Different(t *testing.T) {
	r, _ := diff.Compare("0 * * * *", "30 * * * *")
	s := diff.Summary(r)
	if !strings.Contains(s, "minute") {
		t.Errorf("expected 'minute' in summary, got: %s", s)
	}
}
