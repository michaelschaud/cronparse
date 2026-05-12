package next_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/cronparse/next"
)

func TestFormatText_ContainsExpression(t *testing.T) {
	results := next.ForExpressions([]string{"* * * * *"}, refTime, 2)
	out := next.FormatText(results)
	if !strings.Contains(out, "* * * * *") {
		t.Error("expected expression in text output")
	}
}

func TestFormatText_ContainsRuns(t *testing.T) {
	results := next.ForExpressions([]string{"* * * * *"}, refTime, 2)
	out := next.FormatText(results)
	if !strings.Contains(out, "[1]") || !strings.Contains(out, "[2]") {
		t.Error("expected run indices in text output")
	}
}

func TestFormatText_ErrorShown(t *testing.T) {
	results := next.ForExpressions([]string{"bad"}, refTime, 2)
	out := next.FormatText(results)
	if !strings.Contains(out, "Error") {
		t.Error("expected error label in text output")
	}
}

func TestFormatMarkdown_ContainsTable(t *testing.T) {
	results := next.ForExpressions([]string{"0 * * * *"}, refTime, 3)
	out := next.FormatMarkdown(results)
	if !strings.Contains(out, "|") {
		t.Error("expected markdown table in output")
	}
	if !strings.Contains(out, "Next Run") {
		t.Error("expected 'Next Run' header in markdown")
	}
}

func TestFormatMarkdown_ErrorShown(t *testing.T) {
	results := next.ForExpressions([]string{"nope"}, refTime, 2)
	out := next.FormatMarkdown(results)
	if !strings.Contains(out, "Error") {
		t.Error("expected error in markdown output")
	}
}

func TestFormatMergedText_ContainsHeader(t *testing.T) {
	times := []time.Time{refTime.Add(time.Minute), refTime.Add(2 * time.Minute)}
	out := next.FormatMergedText(times)
	if !strings.Contains(out, "Merged next runs") {
		t.Error("expected header in merged text output")
	}
	if !strings.Contains(out, "[1]") {
		t.Error("expected index in merged text output")
	}
}
