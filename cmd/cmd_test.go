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

	testArg(t, http.MethodPatch, mockRunner.Method, "method")

	if mockRunner.Body != nil {
		testArg(t, "", string(mockRunner.Body), "body")
	}

	testArg(t, "https://localhost:9999/bookings/9999", mockRunner.URL, "URL")

	gotStatus, ok := mockRunner.Query.Params["status"]
	if !ok {
		t.Error("Missing query parameter status")
	}

	testArg(t, "CONFIRMED", gotStatus, "status")

	gotMessage, ok := mockRunner.Query.Params["message"]
	if !ok {
		t.Error("Missing query parameter message")
	}

	testArg(t, "test message", gotMessage, "message")
}

func testArg(t *testing.T, expected, got, argumentName string) {
	t.Helper()
	if expected != got {
		t.Logf("Unexpected %s in command.", argumentName)
		t.Logf("Expected %s", expected)
		t.Logf("Got %s", got)
		t.Fail()
	}
}
