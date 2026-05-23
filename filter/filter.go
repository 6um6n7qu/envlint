package filter

import "github.com/user/envlint/schema"

// Options controls which variables are selected.
type Options struct {
	Groups   []string
	Required *bool
	Prefix   string
}

// Apply returns only the schema variables that match all active filter criteria.
func Apply(vars []schema.VarDefinition, opts Options) []schema.VarDefinition {
	var out []schema.VarDefinition
	for _, v := range vars {
		if opts.Prefix != "" && !hasPrefix(v.Name, opts.Prefix) {
			continue
		}
		if len(opts.Groups) > 0 && !inGroups(v.Group, opts.Groups) {
			continue
		}
		if opts.Required != nil && v.Required != *opts.Required {
			continue
		}
		out = append(out, v)
	}
	return out
}

func hasPrefix(name, prefix string) bool {
	if len(name) < len(prefix) {
		return false
	}
	return name[:len(prefix)] == prefix
}

func inGroups(group string, groups []string) bool {
	for _, g := range groups {
		if g == group {
			return true
		}
	}
	return false
}
