package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/example/cronparse/forecast"
	"github.com/example/cronparse/humanize"
	"github.com/example/cronparse/parser"
)

func main() {
	n := flag.Int("n", 5, "number of next run times to display")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: cronparse [flags] \"<cron expression>\"\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  cronparse \"0 9 * * 1-5\"\n")
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	expr := strings.Join(flag.Args(), " ")

	fields, err := parser.Parse(expr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid cron expression: %v\n", err)
		os.Exit(1)
	}

	description := humanize.Describe(fields)
	fmt.Printf("Expression : %s\n", expr)
	fmt.Printf("Description: %s\n", description)

	now := time.Now()
	runs, err := forecast.NextRuns(expr, now, *n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not compute next runs: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nNext %s runs (from %s):\n", strconv.Itoa(*n), now.Format("2006-01-02 15:04:05"))
	for i, t := range runs {
		fmt.Printf("  %d. %s\n", i+1, t.Format("2006-01-02 15:04:05 (Mon)"))
	}
}
