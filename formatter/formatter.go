package formatter

import (
	"fmt"
	"strings"
)

// Format controls the output format for env variable display.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatExport Format = "export"
	FormatJSON   Format = "json"
)

// Entry represents a single key-value pair to format.
type Entry struct {
	Key   string
	Value string
}

// Render formats a slice of entries according to the given format.
func Render(entries []Entry, format Format) (string, error) {
	switch format {
	case FormatDotenv:
		return renderDotenv(entries), nil
	case FormatExport:
		return renderExport(entries), nil
	case FormatJSON:
		return renderJSON(entries), nil
	default:
		return "", fmt.Errorf("unknown format: %q", format)
	}
}

func renderDotenv(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
	}
	return sb.String()
}

func renderExport(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "export %s=%q\n", e.Key, e.Value)
	}
	return sb.String()
}

func renderJSON(entries []Entry) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, e := range entries {
		comma := ","
		if i == len(entries)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  %q: %q%s\n", e.Key, e.Value, comma)
	}
	sb.WriteString("}\n")
	return sb.String()
}
