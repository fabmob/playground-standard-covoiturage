package test

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
)

// Flags stores validation options
type Flags struct {
	// If true, an Empty response is considered to be an error
	ExpectNonEmpty bool

	// It true, the HTTP response code is tested against expectation
	ExpectedResponseCode int

	// If true, the Booking status retrieved from `Get /bookings` call is tested
	// against expectation
	ExpectedBookingStatus api.BookingStatus

	// If true, the API is supposed to support the booking by deep link use case
	ExpectDeepLinkSupport bool
}

const (
	DefaultFlagExpectNonEmpty        = false
	DefaultFlagExpectDeepLinkSupport = false
	DefaultFlagExpectedResponseCode  = http.StatusOK
	DefaultFlagExpectedBookingStatus = ""
)

// NewFlags return a set of default flags
func NewFlags() Flags {
	return Flags{
		ExpectNonEmpty:        DefaultFlagExpectNonEmpty,
		ExpectDeepLinkSupport: DefaultFlagExpectDeepLinkSupport,
		ExpectedResponseCode:  DefaultFlagExpectedResponseCode,
		ExpectedBookingStatus: DefaultFlagExpectedBookingStatus,
	}
}
