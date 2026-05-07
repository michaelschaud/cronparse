package schedule

import (
	"fmt"
	"time"

	"github.com/yourorg/cronparse/forecast"
	"github.com/yourorg/cronparse/parser"
)

// FrequencySummary describes how often a cron expression fires.
type FrequencySummary struct {
	Expression    string
	RunsPerHour   float64
	RunsPerDay    float64
	AvgInterval   time.Duration
}

// Frequency calculates an approximate frequency summary for the given cron
// expression by sampling the next sampleSize occurrences from base.
func Frequency(expr string, base time.Time, sampleSize int) (*FrequencySummary, error) {
	if sampleSize < 2 {
		sampleSize = 2
	}
	_, err := parser.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("invalid expression: %w", err)
	}

	runs, err := forecast.NextRuns(expr, base, sampleSize)
	if err != nil {
		return nil, err
	}
	if len(runs) < 2 {
		return nil, fmt.Errorf("not enough runs to calculate frequency")
	}

	totalDuration := runs[len(runs)-1].Sub(runs[0])
	avgInterval := totalDuration / time.Duration(len(runs)-1)

	var runsPerHour, runsPerDay float64
	if avgInterval > 0 {
		runsPerHour = float64(time.Hour) / float64(avgInterval)
		runsPerDay = float64(24*time.Hour) / float64(avgInterval)
	}

	return &FrequencySummary{
		Expression:  expr,
		RunsPerHour: runsPerHour,
		RunsPerDay:  runsPerDay,
		AvgInterval: avgInterval,
	}, nil
}
