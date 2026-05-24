// Package profiler analyzes a loaded .env file and produces a summary
// of variable statistics such as counts by group, required vs optional,
// and variables with default values.
package profiler

import "github.com/user/envlint/schema"

// Report holds aggregated statistics about a schema's variables.
type Report struct {
	Total        int
	Required     int
	Optional     int
	WithDefault  int
	WithPattern  int
	WithAllowed  int
	ByGroup      map[string]int
}

// Profile inspects a schema and returns a Report.
func Profile(s *schema.Schema) Report {
	r := Report{
		ByGroup: make(map[string]int),
	}

	for _, v := range s.Vars {
		r.Total++

		if v.Required {
			r.Required++
		} else {
			r.Optional++
		}

		if v.Default != "" {
			r.WithDefault++
		}

		if v.Pattern != "" {
			r.WithPattern++
		}

		if len(v.Allowed) > 0 {
			r.WithAllowed++
		}

		group := v.Group
		if group == "" {
			group = "(ungrouped)"
		}
		r.ByGroup[group]++
	}

	return r
}
