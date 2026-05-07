package explain_test

import (
	"strings"
	"testing"

	"github.com/yourorg/cronparse/explain"
)

func TestFormatText_ContainsExpression(t *testing.T) {
	b, err := explain.Explain("0 12 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := explain.FormatText(b)
	if !strings.Contains(out, "0 12 * * *") {
		t.Errorf("expected expression in text output, got:\n%s", out)
	}
	if !strings.Contains(out, "Minute") {
		t.Errorf("expected 'Minute' field label in output")
	}
	if !strings.Contains(out, "Hour") {
		t.Errorf("expected 'Hour' field label in output")
	}
}

func TestFormatMarkdown_ContainsTable(t *testing.T) {
	b, err := explain.Explain("*/5 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := explain.FormatMarkdown(b)
	if !strings.Contains(out, "| Field |") {
		t.Errorf("expected markdown table header, got:\n%s", out)
	}
	if !strings.Contains(out, "*/5") {
		t.Errorf("expected raw value '*/5' in markdown output")
	}
	if !strings.Contains(out, "**Expression:**") {
		t.Errorf("expected bold Expression label in markdown output")
	}
}

func TestFormatText_AllFieldsPresent(t *testing.T) {
	b, err := explain.Explain("30 9 15 3 5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := explain.FormatText(b)
	expectedFields := []string{"Minute", "Hour", "Day-of-Month", "Month", "Day-of-Week"}
	for _, f := range expectedFields {
		if !strings.Contains(out, f) {
			t.Errorf("expected field %q in text output", f)
		}
	}
}

func TestFormatMarkdown_AllFieldsPresent(t *testing.T) {
	b, err := explain.Explain("30 9 15 3 5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := explain.FormatMarkdown(b)
	expectedFields := []string{"Minute", "Hour", "Day-of-Month", "Month", "Day-of-Week"}
	for _, f := range expectedFields {
		if !strings.Contains(out, f) {
			t.Errorf("expected field %q in markdown output", f)
		}
	}
}
