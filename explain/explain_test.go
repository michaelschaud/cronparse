package explain_test

import (
	"strings"
	"testing"

	"github.com/yourorg/cronparse/explain"
)

func TestExplain_EveryMinute(t *testing.T) {
	b, err := explain.Explain("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Expression != "* * * * *" {
		t.Errorf("expected expression '* * * * *', got %q", b.Expression)
	}
	if len(b.Fields) != 5 {
		t.Fatalf("expected 5 fields, got %d", len(b.Fields))
	}
	if !strings.Contains(b.Fields[0].Meaning, "every") {
		t.Errorf("expected 'every' in minute meaning, got %q", b.Fields[0].Meaning)
	}
}

func TestExplain_SpecificValues(t *testing.T) {
	b, err := explain.Explain("30 9 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Fields[0].Meaning != "30" {
		t.Errorf("expected minute meaning '30', got %q", b.Fields[0].Meaning)
	}
	if b.Fields[1].Meaning != "9" {
		t.Errorf("expected hour meaning '9', got %q", b.Fields[1].Meaning)
	}
	if b.Fields[4].Meaning != "Monday" {
		t.Errorf("expected weekday 'Monday', got %q", b.Fields[4].Meaning)
	}
}

func TestExplain_StepValue(t *testing.T) {
	b, err := explain.Explain("*/15 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(b.Fields[0].Meaning, "15") {
		t.Errorf("expected step '15' in meaning, got %q", b.Fields[0].Meaning)
	}
}

func TestExplain_MonthName(t *testing.T) {
	b, err := explain.Explain("0 0 1 6 *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Fields[3].Meaning != "June" {
		t.Errorf("expected 'June', got %q", b.Fields[3].Meaning)
	}
}

func TestExplain_InvalidExpression(t *testing.T) {
	_, err := explain.Explain("bad expression")
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}
