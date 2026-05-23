package reporter

import (
	"fmt"
	"io"

	"github.com/envlint/envlint/validator"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// Print writes a human-readable validation report to the given writer.
func Print(results []validator.Result, w io.Writer) {
	passed := 0
	failed := 0

	for _, r := range results {
		if r.Valid {
			fmt.Fprintf(w, "%s✔ %s%s\n", colorGreen, r.Key, colorReset)
			passed++
		} else {
			fmt.Fprintf(w, "%s✘ %s: %s%s\n", colorRed, r.Key, r.Message, colorReset)
			failed++
		}
	}

	fmt.Fprintln(w)
	if failed > 0 {
		fmt.Fprintf(w, "%sFailed: %d, Passed: %d%s\n", colorRed, failed, passed, colorReset)
	} else {
		fmt.Fprintf(w, "%sAll %d checks passed%s\n", colorGreen, passed, colorReset)
	}
}

// PrintJSON writes a JSON-formatted validation report to the given writer.
func PrintJSON(results []validator.Result, w io.Writer) {
	fmt.Fprintln(w, "[")
	for i, r := range results {
		comma := ","
		if i == len(results)-1 {
			comma = ""
		}
		valid := "true"
		if !r.Valid {
			valid = "false"
		}
		fmt.Fprintf(w, "  {\"key\": \"%s\", \"valid\": %s, \"message\": \"%s\"}%s\n",
			r.Key, valid, r.Message, comma)
	}
	fmt.Fprintln(w, "]")
}
