package masker

import "flag"

// Options holds CLI flag values for masking configuration.
type Options struct {
	// MaskSecrets controls whether sensitive values are masked in output.
	MaskSecrets bool
	// ExtraPatterns holds additional user-supplied sensitive key patterns.
	ExtraPatterns []string
}

// FlagSet registers masker-related flags onto the given FlagSet and returns
// a pointer to the populated Options struct.
func FlagSet(fs *flag.FlagSet) *Options {
	opts := &Options{}
	fs.BoolVar(&opts.MaskSecrets, "mask-secrets", true, "mask sensitive values in output (default: true)")
	return opts
}

// Build constructs a Masker from the parsed Options. If MaskSecrets is false,
// it returns a Masker with no patterns so nothing is masked.
func (o *Options) Build() *Masker {
	if !o.MaskSecrets {
		return NewWithPatterns([]string{})
	}
	patterns := append([]string{}, DefaultSensitivePatterns...)
	patterns = append(patterns, o.ExtraPatterns...)
	return NewWithPatterns(patterns)
}
