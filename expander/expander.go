// Package expander handles variable interpolation in .env files.
// It resolves references like ${VAR} or $VAR using already-loaded values
// or fallback to OS environment variables.
package expander

import (
	"os"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Expand resolves variable references within the values of the provided env map.
// Resolution order: env map itself (already resolved entries), then OS environment.
// Unresolvable references are left as empty strings.
func Expand(env map[string]string) map[string]string {
	resolved := make(map[string]string, len(env))

	// Copy originals first so self-references can fall back to OS.
	for k, v := range env {
		resolved[k] = v
	}

	for k, v := range resolved {
		resolved[k] = expandValue(v, resolved)
	}

	return resolved
}

// expandValue replaces all variable references in a single value string.
func expandValue(value string, env map[string]string) string {
	return refPattern.ReplaceAllStringFunc(value, func(match string) string {
		name := extractName(match)
		if val, ok := env[name]; ok {
			return val
		}
		if val, ok := os.LookupEnv(name); ok {
			return val
		}
		return ""
	})
}

// extractName pulls the variable name out of a $VAR or ${VAR} reference.
func extractName(ref string) string {
	ref = strings.TrimPrefix(ref, "$")
	ref = strings.TrimPrefix(ref, "{")
	ref = strings.TrimSuffix(ref, "}")
	return ref
}
