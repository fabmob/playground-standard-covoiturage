package service

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
)

var generateTestData bool
var generatedData = NewMockDB()
var commandsFile = strings.Builder{}

// appendDataIfGenerated is used to populate the `generatedData` db if the
// -generate flag is provided
func appendDataIfGenerated(mockDB *MockDB) {
	if generateTestData {
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

	cmdContinuation := " \\\n  "

	cmd += fmt.Sprintf("echo \"%s\"\n", t.Name())
	cmd += "go run main.go test" + cmdContinuation +
		fmt.Sprintf("--method=%s", request.Method) + cmdContinuation +
		fmt.Sprintf("--url=%s", request.URL) + cmdContinuation +
		fmt.Sprintf("--expectStatus=%d", flags.ExpectedStatusCode)

	if body != nil {
		cmd += cmdContinuation + fmt.Sprintf("<<< '%s'", body)
	}

	cmd += "\n\n"
	return cmd
}
