// Package reporter handles formatted output of validation results.
package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/user/envlint/suggester"
	"github.com/user/envlint/validator"
)

// Print writes a human-readable summary of results to stdout.
func Print(results []validator.Result) {
	PrintTo(os.Stdout, results, false)
}

// PrintWithSuggestions writes results including fix hints to stdout.
func PrintWithSuggestions(results []validator.Result) {
	PrintTo(os.Stdout, results, true)
}

// PrintTo writes results to the given writer, optionally including suggestions.
func PrintTo(w io.Writer, results []validator.Result, withSuggestions bool) {
	failed := 0
	for _, r := range results {
		status := "PASS"
		if !r.Valid {
			status = "FAIL"
			failed++
		}
		fmt.Fprintf(w, "[%s] %s: %s\n", status, r.Variable, r.Message)
	}

	if withSuggestions && failed > 0 {
		suggestions := suggester.Suggest(results)
		if len(suggestions) > 0 {
			fmt.Fprintln(w, "\nSuggestions:")
			for _, s := range suggestions {
				fmt.Fprintf(w, "  • %s: %s\n", s.Variable, s.Hint)
			}
		}
	}

	fmt.Fprintf(w, "\nResult: %d passed, %d failed\n", len(results)-failed, failed)
}

// PrintJSON writes results as a JSON array to stdout.
func PrintJSON(results []validator.Result) error {
	return PrintJSONTo(os.Stdout, results)
}

// PrintJSONTo writes results as a JSON array to the given writer.
func PrintJSONTo(w io.Writer, results []validator.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		return fmt.Errorf("failed to encode results as JSON: %w", err)
	}
	return nil
}
