package loader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Env holds the parsed key-value pairs from a .env file.
type Env map[string]string

// Load reads a .env file from the given path and returns the parsed key-value
// pairs. Lines starting with '#' are treated as comments and ignored. Empty
// lines are also skipped. Values may optionally be quoted with single or double
// quotes, which are stripped during parsing.
func Load(path string) (Env, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loader: cannot open file %q: %w", path, err)
	}
	defer f.Close()

	env := make(Env)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("loader: invalid syntax at line %d: %q", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("loader: empty key at line %d", lineNum)
		}

		value = stripQuotes(value)
		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("loader: error reading file: %w", err)
	}

	return env, nil
}

// stripQuotes removes surrounding single or double quotes from a value if
// both the opening and closing quote characters match.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
