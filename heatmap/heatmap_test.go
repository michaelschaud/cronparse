package heatmap_test

import (
	"testing"
	"time"

	"github.com/cronparse/heatmap"
)

// Monday 2024-01-01 00:00:00 UTC
var fixedFrom = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestBuild_EveryMinuteOneHour(t *testing.T) {
	to := fixedFrom.Add(time.Hour)
	m, err := heatmap.Build("* * * * *", fixedFrom, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Expression != "* * * * *" {
		t.Errorf("expression mismatch")
	}
	// Monday = dow 1, hour 0 should have 60 fires.
	if m.Grid[1][0] != 60 {
		t.Errorf("expected 60 fires in Mon/00, got %d", m.Grid[1][0])
	}
}

func TestBuild_HourlyTwoDays(t *testing.T) {
	to := fixedFrom.Add(48 * time.Hour)
	m, err := heatmap.Build("0 * * * *", fixedFrom, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Every hour fires once; over 48 hours we expect 48 total fires.
	total := 0
	for _, c := range m.Cells {
		total += c.Count
	}
	if total != 48 {
		t.Errorf("expected 48 total fires, got %d", total)
	}
}

func TestBuild_InvalidExpression(t *testing.T) {
	_, err := heatmap.Build("bad expr", fixedFrom, fixedFrom.Add(time.Hour))
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestBuild_ToBeforeFrom(t *testing.T) {
	_, err := heatmap.Build("* * * * *", fixedFrom, fixedFrom.Add(-time.Minute))
	if err == nil {
		t.Fatal("expected error when to <= from")
	}
}

func TestPeak_ReturnsHighestCount(t *testing.T) {
	to := fixedFrom.Add(time.Hour)
	m, err := heatmap.Build("* * * * *", fixedFrom, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p := m.Peak()
	if p.Count == 0 {
		t.Fatal("expected non-zero peak")
	}
	// All fires are in hour 0 of the window.
	if p.Hour != 0 {
		t.Errorf("expected peak hour 0, got %d", p.Hour)
	}
}

func TestPeak_EmptyMap(t *testing.T) {
	m := &heatmap.Map{}
	p := m.Peak()
	if p.Count != 0 {
		t.Errorf("expected zero peak for empty map")
	}
}
