package test

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
)

// Flags stores validation options
type Flags struct {
	// If true, an Empty response is considered to be an error
	DisallowEmpty bool

	// If true, the API is supposed to support the booking by deep link use case
	SupportDeepLink bool

	ExpectedStatusCode int

	ExpectedBookingStatus api.BookingStatus
}

const (
	DefaultDisallowEmptyFlag   = false
	DefaultSupportDeepLinkFlag = false
)

// NewFlags return a set of default flags
func NewFlags() Flags {
	return Flags{
		DisallowEmpty:         false,
		SupportDeepLink:       false,
		ExpectedStatusCode:    http.StatusOK,
		ExpectedBookingStatus: "",
	}
}
