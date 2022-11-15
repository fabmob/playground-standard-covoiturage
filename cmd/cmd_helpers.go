package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
	"github.com/stoewer/go-strcase"
)

// makeEndpointCommand creates a cobra command skeletton for a given endpoint
func makeEndpointCommand(endpoint test.Endpoint) *cobra.Command {
	pathNoLeadingSlash := strings.TrimPrefix(endpoint.Path, "/")
	return &cobra.Command{
		Use:   strcase.LowerCamelCase(pathNoLeadingSlash),
		Short: cmdDescription(endpoint),
		Long:  cmdDescription(endpoint),
	}
}

/////////////////////////////////////////////////////

var (
	checkRequiredBookingID          = checkRequiredString("bookingId")
	checkRequiredStatus             = checkRequiredString("status")
	checkRequiredServer             = checkRequiredString("server")
	checkRequiredDepartureLat       = checkRequiredString("departureLat")
	checkRequiredDepartureLng       = checkRequiredString("departureLng")
	checkRequiredArrivalLat         = checkRequiredString("arrivalLat")
	checkRequiredArrivalLng         = checkRequiredString("arrivalLng")
	checkRequiredDepartureDate      = checkRequiredString("departureDate")
	checkRequiredDepartureTimeOfDay = checkRequiredString("departureTimeOfDay")
)

// checkRequiredString is a partial application that helps creating testing
// functions for non-empty string flags
func checkRequiredString(description string) func(string) error {
	return func(obj string) error {
		if obj == "" {
			return fmt.Errorf("missing required --%s information", description)
		}

		return nil
	}
}

// A short command description for a given endpoint
func cmdDescription(endpoint test.Endpoint) string {
	return fmt.Sprintf("Test the %s endpoint", endpoint)
}

// readBodyFromStdin reads stdin stream until it is closed, and returns its
// content. The function returns an error if it is not closed before `timeout`, or if an error occurs while
// reading.
func readBodyFromStdin(cmd *cobra.Command, timeout time.Duration) ([]byte,
	error) {

	var (
		stdinStreamReader = cmd.InOrStdin()
		stdinChannel      = make(chan []byte, 1)
		errChannel        = make(chan error, 1)
	)

	go func() {
		b, err := io.ReadAll(stdinStreamReader)
		if err != nil {
			errChannel <- err
		} else {
			stdinChannel <- b
		}
	}()

	select {
	case <-time.After(timeout):
		return nil, errors.New("body is required but missing")

	case err := <-errChannel:
		return nil, err

	case body := <-stdinChannel:
		return body, nil
	}
}

// anyError returns first non-nil error (or nil if none exists)
func anyError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
