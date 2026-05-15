package sparse_test

import (
	"testing"
	"time"

	"github.com/cronparse/sparse"
)

func fixedFrom() time.Time {
	return time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
}

func TestFind_EveryMinuteNoSparseGaps(t *testing.T) {
	from := fixedFrom()
	to := from.Add(2 * time.Hour)
	res := sparse.Find("* * * * *", from, to, 5*time.Minute)
	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if len(res.Gaps) != 0 {
		t.Errorf("expected no sparse gaps for every-minute expr, got %d", len(res.Gaps))
	}
}

func TestFind_HourlyHasSparseGaps(t *testing.T) {
	from := fixedFrom()
	to := from.Add(6 * time.Hour)
	res := sparse.Find("0 * * * *", from, to, 30*time.Minute)
	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if len(res.Gaps) == 0 {
		t.Error("expected sparse gaps for hourly expr with 30m threshold")
	}
	for _, g := range res.Gaps {
		if g.Duration <= 30*time.Minute {
			t.Errorf("gap duration %v should exceed minGap", g.Duration)
		}
	}
}

func TestFind_LongestIsMaxGap(t *testing.T) {
	from := fixedFrom()
	to := from.Add(6 * time.Hour)
	res := sparse.Find("0 * * * *", from, to, 30*time.Minute)
	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	for _, g := range res.Gaps {
		if g.Duration > res.Longest {
			t.Errorf("gap %v exceeds reported Longest %v", g.Duration, res.Longest)
		}
	}
}

func TestFind_InvalidExpression(t *testing.T) {
	from := fixedFrom()
	to := from.Add(time.Hour)
	res := sparse.Find("invalid", from, to, time.Minute)
	if res.Err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestFind_ToBeforeFrom(t *testing.T) {
	from := fixedFrom()
	to := from.Add(-time.Hour)
	res := sparse.Find("* * * * *", from, to, time.Minute)
	if res.Err == nil {
		t.Error("expected error when to is before from")
	}
}

func TestFind_GapFieldsPopulated(t *testing.T) {
	from := fixedFrom()
	to := from.Add(4 * time.Hour)
	res := sparse.Find("0 * * * *", from, to, time.Minute)
	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	for _, g := range res.Gaps {
		if g.To.Before(g.From) || g.To.Equal(g.From) {
			t.Errorf("gap To %v is not after From %v", g.To, g.From)
		}
		if g.Duration != g.To.Sub(g.From) {
			t.Errorf("gap Duration %v does not match To-From %v", g.Duration, g.To.Sub(g.From))
		}
	}
}
