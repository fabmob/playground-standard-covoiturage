package cmd

import (
	"net/http"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
)

func TestPatchBookingsCmd(t *testing.T) {

	var (
		server         = "https://localhost:9999"
		bookingID      = "9999"
		status         = "CONFIRMED"
		message        = "test message"
		expectedURL    = "https://localhost:9999/bookings/9999"
		expectedMethod = http.MethodPatch
	)

	mockRunner := test.NewMockRunner()
	err := patchBookingsRun(
		mockRunner,
		server,
		bookingID,
		status,
		message,
		test.NewFlags(),
	)
	panicIf(err)

	// Test Assertions
	testStringArg(t, expectedMethod, mockRunner.Method, "method")

	// Nil or empty body
	if mockRunner.Body != nil {
		testStringArg(t, string(mockRunner.Body), "", "body")
	}

	testStringArg(t, expectedURL, mockRunner.URL, "URL")

	gotStatus, ok := mockRunner.Query.Params["status"]
	if !ok {
		t.Error("Missing query parameter status")
	}

	testStringArg(t, gotStatus, status, "status")

	gotMessage, ok := mockRunner.Query.Params["message"]
	if !ok {
		t.Error("Missing query parameter message")
	}

	testStringArg(t, gotMessage, message, "message")
}

func TestPostMessages(t *testing.T) {

	var (
		server             = "https://localhost:9999"
		expectedBody       = "body"
		bodyBytes          = []byte(expectedBody)
		expectedMethod     = http.MethodPost
		expectedStatusCode = http.StatusCreated
	)

	mockRunner := test.NewMockRunner()
	err := getMessagesRun(mockRunner, server, bodyBytes)
	panicIf(err)

	// Test Assertions

	testStringArg(t, mockRunner.Method, expectedMethod, "method")

	testStringArg(t, mockRunner.URL, "https://localhost:9999/messages", "URL")

	if mockRunner.Body == nil {
		t.Error("Missing required body")
	}

	testStringArg(t, string(mockRunner.Body), expectedBody, "body")

	if mockRunner.Flags.ExpectedStatusCode != expectedStatusCode {
		t.Error("Wrong default status code")
	}

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

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
