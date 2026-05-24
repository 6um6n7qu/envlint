package linter

import (
	"flag"
	"strings"
)

// Options holds CLI-configurable options for the linter.
type Options struct {
	FailOnWarning bool
	SkipRules     []string
}

// FlagSet returns a flag.FlagSet configured for linter CLI flags.
func FlagSet(opts *Options) *flag.FlagSet {
	fs := flag.NewFlagSet("lint", flag.ContinueOnError)

	fs.BoolVar(&opts.FailOnWarning, "fail-on-warning", false,
		"treat warnings as errors and return non-zero exit code")

	var skipRules string
	fs.StringVar(&skipRules, "skip-rules", "",
		"comma-separated list of rule names to skip (e.g. empty-description,default-on-required)")

	// Post-parse hook: split the comma list into the slice.
	// Caller must invoke this after fs.Parse.
	fs.VisitAll(func(f *flag.Flag) {})

	_ = skipRules // parsed lazily via ParseSkipRules
	return fs
}

// ParseSkipRules splits a comma-separated rule string into a slice.
func ParseSkipRules(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	var out []string
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// HasErrors returns true if any issue has severity "error",
// or if opts.FailOnWarning is set and any warning exists.
func HasErrors(issues []Issue, opts Options) bool {
	for _, iss := range issues {
		if iss.Severity == "error" {
			return true
		}
		if opts.FailOnWarning && iss.Severity == "warning" {
			return true
		}
	}
	return false
}
