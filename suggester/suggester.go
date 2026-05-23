// Package suggester provides hints and fix suggestions for validation errors.
package suggester

import (
	"fmt"
	"strings"

	"github.com/user/envlint/validator"
)

// Suggestion holds a human-readable hint for a validation issue.
type Suggestion struct {
	Variable string
	Hint     string
}

// Suggest generates fix suggestions based on a slice of validation results.
func Suggest(results []validator.Result) []Suggestion {
	var suggestions []Suggestion

	for _, r := range results {
		if r.Valid {
			continue
		}
		s := Suggestion{
			Variable: r.Variable,
			Hint:     buildHint(r),
		}
		suggestions = append(suggestions, s)
	}

	return suggestions
}

func buildHint(r validator.Result) string {
	msg := strings.ToLower(r.Message)

	switch {
	case strings.Contains(msg, "missing") && strings.Contains(msg, "required"):
		return fmt.Sprintf("Add '%s' to your .env file. It is required for the application to run.", r.Variable)
	case strings.Contains(msg, "pattern"):
		return fmt.Sprintf("Check the format of '%s'. Ensure it matches the expected pattern defined in the schema.", r.Variable)
	case strings.Contains(msg, "allowed"):
		return fmt.Sprintf("'%s' has an invalid value. Refer to the schema for the list of allowed values.", r.Variable)
	default:
		return fmt.Sprintf("Review the value of '%s' against the schema definition.", r.Variable)
	}
}
