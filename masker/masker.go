// Package masker provides utilities for masking sensitive environment variable values
// in output to prevent accidental secret exposure in logs or reports.
package masker

import "strings"

// DefaultSensitivePatterns contains common substrings that indicate a variable
// holds a sensitive value and should be masked.
var DefaultSensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
}

const maskValue = "***"

// Masker holds configuration for masking sensitive values.
type Masker struct {
	Patterns []string
}

// New returns a Masker using the default sensitive patterns.
func New() *Masker {
	return &Masker{Patterns: DefaultSensitivePatterns}
}

// NewWithPatterns returns a Masker using the provided patterns.
func NewWithPatterns(patterns []string) *Masker {
	return &Masker{Patterns: patterns}
}

// IsSensitive reports whether the given key matches any sensitive pattern.
func (m *Masker) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range m.Patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// Mask returns the masked representation of a value if the key is sensitive,
// otherwise returns the original value.
func (m *Masker) Mask(key, value string) string {
	if m.IsSensitive(key) {
		return maskValue
	}
	return value
}

// MaskAll applies masking to an entire map of key-value pairs, returning a new
// map with sensitive values replaced.
func (m *Masker) MaskAll(vars map[string]string) map[string]string {
	result := make(map[string]string, len(vars))
	for k, v := range vars {
		result[k] = m.Mask(k, v)
	}
	return result
}
