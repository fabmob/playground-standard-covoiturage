package test

type Flags struct {
	DisallowEmpty   bool
	SupportDeepLink bool
}

// NewFlags return a set of default flags
func NewFlags() Flags {
	return Flags{
		DisallowEmpty:   false,
		SupportDeepLink: false,
	}
}
