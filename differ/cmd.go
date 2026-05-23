package differ

import "flag"

// Options holds CLI flags for the differ subcommand.
type Options struct {
	BaseFile  string
	OtherFile string
	Quiet     bool
}

// FlagSet returns a *flag.FlagSet pre-configured for the differ subcommand.
func FlagSet(opts *Options) *flag.FlagSet {
	fs := flag.NewFlagSet("diff", flag.ContinueOnError)
	fs.StringVar(&opts.BaseFile, "base", ".env", "base .env file to compare from")
	fs.StringVar(&opts.OtherFile, "other", ".env.example", "other .env file to compare against")
	fs.BoolVar(&opts.Quiet, "quiet", false, "suppress output; exit code only")
	return fs
}
