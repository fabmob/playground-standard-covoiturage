package service

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
)

//go:generate bash -c "go test -generate > /dev/null"

var generateTestData bool
var generatedData = NewMockDB()
var commandsFile = strings.Builder{}
var serverEnvVar = "SERVER"
var authEnvVar = "API_TOKEN"

func init() {
	fmt.Fprintln(&commandsFile, "#!/usr/bin/env bash")
	fmt.Fprint(&commandsFile, "# Generated programmatically - DO NOT EDIT\n\n")
	fmt.Fprintf(&commandsFile, "export %s=\"%s\"\n", serverEnvVar, localServer)
	fmt.Fprintf(&commandsFile, "export %s=\"\"\n\n", authEnvVar)

}

// Data needs to be appended once for each test, so we keep track if data has
// already been appended for a given test (with test.Name() as key)
var hasAlreadyAppended = map[string]bool{}

// appendDataIfGenerated is used to populate the `generatedData` db if the
// -generate flag is provided
func appendDataIfGenerated(t *testing.T, mockDB *MockDB) {
	if _, ok := hasAlreadyAppended[t.Name()]; generateTestData && !ok {
		hasAlreadyAppended[t.Name()] = true
		appendData(mockDB, generatedData)
	}
}

// appendCmdIfGenerated is used to populate the `commands` string, if
// -generate flag is provided
func appendCmdIfGenerated(t *testing.T, request *http.Request, flags test.Flags, body []byte) {
	if generateTestData {
		fmt.Fprint(
			&commandsFile,
			GenerateCommandStr(t, request, flags, body),
		)
	}
}

func appendData(from *MockDB, to *MockDB) {
	to.DriverJourneys = append(
		to.GetDriverJourneys(),
		from.GetDriverJourneys()...,
	)

	to.PassengerJourneys = append(
		to.GetPassengerJourneys(),
		from.GetPassengerJourneys()...,
	)

	to.Users = append(
		to.GetUsers(),
		from.GetUsers()...,
	)

	for _, booking := range from.GetBookings() {
		err := to.AddBooking(*booking)
		panicIf(err)
	}
}

// GenerateCommandStr generates a string with the command that should be run
// to test the request with given flags and body.
//
// Used to transform golang tests into command line tests.
func GenerateCommandStr(t *testing.T, request *http.Request, flags test.Flags, body []byte) string {
	var cmd string

	urlWithEnvVar := fmt.Sprintf("$%s%s", serverEnvVar,
		strings.TrimPrefix(request.URL.String(), localServer))

	cmdContinuation := " \\\n  "

	cmd += fmt.Sprintf("echo \"%s\"\n", t.Name())
	cmd += "go run main.go test" + cmdContinuation +
		fmt.Sprintf("--method=%s", request.Method) + cmdContinuation +
		fmt.Sprintf("--url=\"%s\"", urlWithEnvVar) + cmdContinuation +
		fmt.Sprintf("--expectResponseCode=%d", flags.ExpectedStatusCode) +
		cmdContinuation +
		fmt.Sprintf("--auth=\"$%s\"", authEnvVar)

	if flags.ExpectedBookingStatus != "" {
		cmd += cmdContinuation + fmt.Sprintf("--expectBookingStatus=%s", flags.ExpectedBookingStatus)
	}

	if body != nil {
		cmd += cmdContinuation + fmt.Sprintf("<<< '%s'", body)
	}

	cmd += "\n\n"
	return cmd
}
