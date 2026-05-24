package profiler

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// PrintTo writes a human-readable profile summary to w.
func PrintTo(w io.Writer, r Report) {
	fmt.Fprintf(w, "Schema Profile\n")
	fmt.Fprintf(w, "  Total variables : %d\n", r.Total)
	fmt.Fprintf(w, "  Required        : %d\n", r.Required)
	fmt.Fprintf(w, "  Optional        : %d\n", r.Optional)
	fmt.Fprintf(w, "  With default    : %d\n", r.WithDefault)
	fmt.Fprintf(w, "  With pattern    : %d\n", r.WithPattern)
	fmt.Fprintf(w, "  With allowed    : %d\n", r.WithAllowed)

	if len(r.ByGroup) == 0 {
		return
	}

	fmt.Fprintf(w, "  By group:\n")
	groups := make([]string, 0, len(r.ByGroup))
	for g := range r.ByGroup {
		groups = append(groups, g)
	}
	sort.Strings(groups)
	for _, g := range groups {
		fmt.Fprintf(w, "    %-20s %d\n", g, r.ByGroup[g])
	}
}

// PrintJSONTo writes the profile as JSON to w.
func PrintJSONTo(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
