package linter

import (
	"fmt"

	"github.com/your-org/envlint/schema"
	"github.com/your-org/envlint/validator"
)

// Rule represents a single lint rule that can be applied to a schema variable.
type Rule func(v schema.Var) []Issue

// Issue represents a lint problem found in the schema definition itself.
type Issue struct {
	Var     string
	Message string
	Severity string // "error" or "warning"
}

// Lint runs all built-in lint rules against the schema and returns any issues found.
func Lint(s *schema.Schema) []Issue {
	rules := []Rule{
		rulePatternAndAllowedValues,
		ruleEmptyDescription,
		ruleDefaultOnRequired,
	}

	var issues []Issue
	for _, v := range s.Vars {
		for _, rule := range rules {
			issues = append(issues, rule(v)...)
		}
	}
	return issues
}

// LintAndValidate runs schema linting and then validates the provided env map.
// Returns lint issues and validation results separately.
func LintAndValidate(s *schema.Schema, env map[string]string) ([]Issue, []validator.Result) {
	issues := Lint(s)
	results := validator.Validate(s, env)
	return issues, results
}

// rulePatternAndAllowedValues warns when both pattern and allowed_values are set.
func rulePatternAndAllowedValues(v schema.Var) []Issue {
	if v.Pattern != "" && len(v.AllowedValues) > 0 {
		return []Issue{{
			Var:      v.Name,
			Message:  "both 'pattern' and 'allowed_values' are set; 'allowed_values' takes precedence",
			Severity: "warning",
		}}
	}
	return nil
}

// ruleEmptyDescription warns when a required variable has no description.
func ruleEmptyDescription(v schema.Var) []Issue {
	if v.Required && v.Description == "" {
		return []Issue{{
			Var:      v.Name,
			Message:  fmt.Sprintf("required variable %q has no description", v.Name),
			Severity: "warning",
		}}
	}
	return nil
}

// ruleDefaultOnRequired errors when a required variable also has a default value.
func ruleDefaultOnRequired(v schema.Var) []Issue {
	if v.Required && v.Default != "" {
		return []Issue{{
			Var:      v.Name,
			Message:  fmt.Sprintf("required variable %q should not have a default value", v.Name),
			Severity: "error",
		}}
	}
	return nil
}
