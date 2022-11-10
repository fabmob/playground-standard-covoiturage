package test

import "net/http"

// Flags stores validation options
type Flags struct {
	// If true, an Empty response is considered to be an error
	DisallowEmpty bool

	// If true, the API is supposed to support the booking by deep link use case
	SupportDeepLink bool

	ExpectedStatusCode int
}

const (
	DefaultDisallowEmptyFlag   = false
	DefaultSupportDeepLinkFlag = false
)

// NewFlags return a set of default flags
func NewFlags() Flags {
	return Flags{
		DisallowEmpty:      false,
		SupportDeepLink:    false,
		ExpectedStatusCode: http.StatusOK,
	}
}
