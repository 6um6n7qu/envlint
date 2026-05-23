// Package filter provides utilities for narrowing schema variable lists
// based on prefix, group membership, or required status.
//
// It is used by the CLI to support targeted linting, e.g.:
//
//	envlint --group database --required
package filter

import (
	"flag"
	"strings"
)

// FlagSet registers filter-related CLI flags and returns a function
// that builds an Options from the parsed flag values.
func FlagSet(fs *flag.FlagSet) func() Options {
	prefix := fs.String("prefix", "", "only check vars with this name prefix")
	groups := fs.String("group", "", "comma-separated list of groups to include")
	requiredOnly := fs.Bool("required", false, "only check required variables")
	optionalOnly := fs.Bool("optional", false, "only check optional variables")

	return func() Options {
		opts := Options{
			Prefix: *prefix,
		}
		if *groups != "" {
			for _, g := range strings.Split(*groups, ",") {
				g = strings.TrimSpace(g)
				if g != "" {
					opts.Groups = append(opts.Groups, g)
				}
			}
		}
		switch {
		case *requiredOnly:
			v := true
			opts.Required = &v
		case *optionalOnly:
			v := false
			opts.Required = &v
		}
		return opts
	}
}
