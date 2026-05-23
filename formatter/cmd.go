package formatter

import (
	"flag"
	"fmt"
)

// Options holds CLI-parsed formatter settings.
type Options struct {
	Format Format
}

// FlagSet registers formatter flags onto the given FlagSet and returns
// a function that builds an Options from the parsed flags.
func FlagSet(fs *flag.FlagSet) func() (*Options, error) {
	format := fs.String("format", string(FormatDotenv),
		`output format: dotenv | export | json`)

	return func() (*Options, error) {
		f := Format(*format)
		switch f {
		case FormatDotenv, FormatExport, FormatJSON:
			// valid
		default:
			return nil, fmt.Errorf("invalid format %q: must be dotenv, export, or json", f)
		}
		return &Options{Format: f}, nil
	}
}
