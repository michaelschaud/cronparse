package tag_test

import (
	"strings"
	"testing"

	"github.com/yourorg/cronparse/tag"
)

func buildCatalog(t *testing.T) *tag.Catalog {
	t.Helper()
	c := tag.NewCatalog()
	_ = c.Add("deployments", "nightly build", "0 2 * * *")
	_ = c.Add("deployments", "weekly release", "0 10 * * 1")
	_ = c.Add("monitoring", "health check", "*/5 * * * *")
	return c
}

func TestFormatText_ContainsGroupName(t *testing.T) {
	c := buildCatalog(t)
	out := tag.FormatText(c)
	if !strings.Contains(out, "[deployments]") {
		t.Errorf("expected '[deployments]' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "[monitoring]") {
		t.Errorf("expected '[monitoring]' in output, got:\n%s", out)
	}
}

func TestFormatText_ContainsLabelsAndExpressions(t *testing.T) {
	c := buildCatalog(t)
	out := tag.FormatText(c)
	if !strings.Contains(out, "nightly build") {
		t.Errorf("expected label 'nightly build' in output")
	}
	if !strings.Contains(out, "0 2 * * *") {
		t.Errorf("expected expression '0 2 * * *' in output")
	}
}

func TestFormatText_EmptyCatalog(t *testing.T) {
	c := tag.NewCatalog()
	out := tag.FormatText(c)
	if !strings.Contains(out, "(no entries)") {
		t.Errorf("expected '(no entries)' for empty catalog, got: %s", out)
	}
}

func TestFormatMarkdown_ContainsTable(t *testing.T) {
	c := buildCatalog(t)
	out := tag.FormatMarkdown(c)
	if !strings.Contains(out, "| Label | Expression |") {
		t.Errorf("expected markdown table header in output")
	}
}

func TestFormatMarkdown_ContainsGroupHeading(t *testing.T) {
	c := buildCatalog(t)
	out := tag.FormatMarkdown(c)
	if !strings.Contains(out, "## deployments") {
		t.Errorf("expected '## deployments' heading in markdown output")
	}
}

func TestFormatMarkdown_EmptyCatalog(t *testing.T) {
	c := tag.NewCatalog()
	out := tag.FormatMarkdown(c)
	if !strings.Contains(out, "_(no entries)_") {
		t.Errorf("expected empty marker in markdown output, got: %s", out)
	}
}
