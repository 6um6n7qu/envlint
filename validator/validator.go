package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/envlint/schema"
)

// Result holds the outcome of validating a single variable.
type Result struct {
	Key     string
	Passed  bool
	Message string
}

// Report aggregates all validation results.
type Report struct {
	Results []Result
}

// HasErrors returns true if any result failed.
func (r *Report) HasErrors() bool {
	for _, res := range r.Results {
		if !res.Passed {
			return true
		}
	}
	return false
}

// Errors returns only the failed results.
func (r *Report) Errors() []Result {
	var errs []Result
	for _, res := range r.Results {
		if !res.Passed {
			errs = append(errs, res)
		}
	}
	return errs
}

// Validate checks the provided env map against the schema.
func Validate(s *schema.Schema, env map[string]string) *Report {
	report := &Report{}

	for _, v := range s.Vars {
		val, exists := env[v.Name]

		if !exists || strings.TrimSpace(val) == "" {
			if v.Required {
				report.Results = append(report.Results, Result{
					Key:     v.Name,
					Passed:  false,
					Message: fmt.Sprintf("%s is required but missing or empty", v.Name),
				})
			}
			continue
		}

		if v.Pattern != "" {
			re, err := regexp.Compile(v.Pattern)
			if err != nil {
				report.Results = append(report.Results, Result{
					Key:     v.Name,
					Passed:  false,
					Message: fmt.Sprintf("%s has invalid pattern in schema: %v", v.Name, err),
				})
				continue
			}
			if !re.MatchString(val) {
				report.Results = append(report.Results, Result{
					Key:     v.Name,
					Passed:  false,
					Message: fmt.Sprintf("%s value %q does not match pattern %q", v.Name, val, v.Pattern),
				})
				continue
			}
		}

		report.Results = append(report.Results, Result{
			Key:    v.Name,
			Passed: true,
		})
	}

	return report
}
