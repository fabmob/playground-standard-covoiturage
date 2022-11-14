package cmd

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var (
	checkRequiredBookingID = checkRequiredString("bookingId")
	checkRequiredStatus    = checkRequiredString("status")
	checkRequiredServer    = checkRequiredString("server")
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

func cmdDescription(endpoint test.Endpoint) string {
	return fmt.Sprintf("Test the %s endpoint", endpoint)
}
