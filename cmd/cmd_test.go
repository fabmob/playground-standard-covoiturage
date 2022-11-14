package cmd

import (
	"net/http"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
)

func TestPatchBookingsCmd(t *testing.T) {
	mockRunner := test.NewMockRunner()

	var (
		server    = "https://localhost:9999"
		bookingID = "9999"
		status    = "CONFIRMED"
		message   = "test message"
	)

	err := patchBookingsRun(mockRunner, server, bookingID, status, message, test.NewFlags())
	if err != nil {
		panic(err)
	}

	testStringArg(t, http.MethodPatch, mockRunner.Method, "method")

	if mockRunner.Body != nil {
		testStringArg(t, string(mockRunner.Body), "", "body")
	}

	testStringArg(t, "https://localhost:9999/bookings/9999", mockRunner.URL, "URL")

	gotStatus, ok := mockRunner.Query.Params["status"]
	if !ok {
		t.Error("Missing query parameter status")
	}

	testStringArg(t, gotStatus, "CONFIRMED", "status")

	gotMessage, ok := mockRunner.Query.Params["message"]
	if !ok {
		t.Error("Missing query parameter message")
	}

	testStringArg(t, gotMessage, "test message", "message")
}

func TestGetMessages(t *testing.T) {
	mockRunner := test.NewMockRunner()

	var (
		server = "https://localhost:9999"
	)

	err := getMessagesRun(mockRunner, server)
	if err != nil {
		panic(err)
	}

	testStringArg(t, mockRunner.URL, "https://localhost:9999", "URL")

}

func testStringArg(t *testing.T, got, expected, argumentName string) {
	t.Helper()
	if expected != got {
		t.Logf("Unexpected %s in command.", argumentName)
		t.Logf("Expected %s", expected)
		t.Logf("Got %s", got)
		t.Fail()
	}
}
