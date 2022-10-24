package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-cli/client"
)

func TestAssertionResult_String(t *testing.T) {
	endpointPath := "/endpoint_path"
	endpointMethod := http.MethodGet
	assertStr := "test assertion"
	errorDescription := "Error description"

	makeAssertionResult := func(err error) AssertionResult {
		return NewAssertionResult(err, endpointPath, endpointMethod, assertStr)
	}
	shouldContain := func(t *testing.T, a AssertionResult, str string) {
		t.Helper()
		if !strings.Contains(a.String(), str) {
			t.Logf("Assertion string : %s", a.String())
			t.Error("Assertion string does not contain " + str)
		}
	}

	testCases := []struct {
		name string
		err  error
	}{
		{
			"Assertion without error",
			nil,
		},
		{
			"Assertion with error",
			errors.New(errorDescription),
		},
	}

	for _, tc := range testCases {
		t.Run("Assertion with error", func(t *testing.T) {
			a := makeAssertionResult(tc.err)
			shouldContain(t, a, endpointMethod)
			shouldContain(t, a, endpointPath)
			shouldContain(t, a, assertStr)
			if tc.err != nil {
				shouldContain(t, a, errorDescription)
			}
		})
	}
}

func TestExpectStatusCode(t *testing.T) {
	testCases := []struct {
		response         *http.Response
		testedStatusCode int
		expectNilError   bool
	}{
		{
			mockStatusResponse(http.StatusOK),
			http.StatusOK,
			true,
		},
		{
			mockStatusResponse(http.StatusTooManyRequests),
			http.StatusTooManyRequests,
			true,
		},
		{
			mockStatusResponse(http.StatusTooManyRequests),
			http.StatusOK,
			false,
		},
		{
			mockStatusResponse(http.StatusInternalServerError),
			http.StatusNotFound,
			false,
		},
	}

	for _, tc := range testCases {
		assertion := assertStatusCode{tc.response, tc.testedStatusCode}
		assertionError := runSingleAssertion(t, assertion)
		if (assertionError == nil) != tc.expectNilError {
			t.Logf("Response status code: %d", tc.response.StatusCode)
			t.Logf("Tested status code: %d", tc.testedStatusCode)
			t.Logf("`expectStatusCode` expected to raise error: %t", !tc.expectNilError)
			t.Error("`expectStatusCode` has not expected behavior")
		}
	}

}

func TestExpectHeaders(t *testing.T) {

	headerContentTypeJSON := http.Header{
		"Content-Type": {"application/json"},
	}
	headerContentTypeJSONWithCharset := http.Header{
		"Content-Type": {"application/json; charset=UTF-8"},
	}
	headerContentTypeForm := http.Header{
		"Content-Type": {"multipart/form-data"},
	}

	testCases := []struct {
		name           string
		header         http.Header
		testKey        string
		testValue      string
		expectNilError bool
	}{
		{
			"No Content-Type header",
			make(http.Header),
			"Content-Type",
			"application/json",
			false,
		},
		{
			"json Content-Type header",
			headerContentTypeJSON,
			"Content-Type",
			"application/json",
			true,
		},
		{
			"json Content-Type header with charset",
			headerContentTypeJSONWithCharset,
			"Content-Type",
			"application/json",
			true,
		},
		{
			"json Content-Type header",
			headerContentTypeJSON,
			"Server",
			"application/json",
			false,
		},
		{
			"wrong Content-Type header",
			headerContentTypeForm,
			"Content-Type",
			"application/json",
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := mockResponse(http.StatusOK, "", tc.header)
			assertion := assertHeaderContains{r, tc.testKey, tc.testValue}
			assertionError := runSingleAssertion(t, assertion)
			if (assertionError == nil) != tc.expectNilError {
				t.Logf("Headers: %v", tc.header)
				t.Logf("Key/value under test: \"%s:%s\"", tc.testKey, tc.testValue)
				t.Logf("AssertHeader expected to raise error: %t", !tc.expectNilError)
				t.Error("AssertHeader has not expected behavior")
			}
		})
	}
}

func TestExpectDriverJourneysFormat(t *testing.T) {

	marshalDriverJourneys := func(dj []client.DriverJourney) string {
		bodyBytes, _ := json.Marshal(dj)
		return string(bodyBytes)
	}

	emptyDriverJourneysBody := marshalDriverJourneys([]client.DriverJourney{})
	singleDriverJourneyBody := marshalDriverJourneys([]client.DriverJourney{{Type: "DYNAMIC"}})
	notAllowedByEnum := marshalDriverJourneys([]client.DriverJourney{{Type: "Not allowed"}})

	missingProp := `[
  {
    "duration": 0,
    "operator": "",
    "passengerDropLat": 0,
    "passengerDropLng": 0,
    "passengerPickupDate": 0,
    "passengerPickupLat": 0,
    "passengerPickupLng": 0,
    "type": ""
  }
]
`

	jsonContentTypeHeader := http.Header{"Content-Type": []string{"application/json"}}

	testCases := []struct {
		name           string
		body           string
		header         http.Header
		expectNilError bool
	}{
		{
			"Not JSON",
			"Hello, world!",
			jsonContentTypeHeader,
			false,
		},
		{
			"Empty []DriverJourney JSON",
			emptyDriverJourneysBody,
			jsonContentTypeHeader,
			true,
		},
		{
			"Non-empty []DriverJourney JSON",
			singleDriverJourneyBody,
			jsonContentTypeHeader,
			true,
		},
		{
			"Other content type",
			"Hello, world!",
			http.Header{"Content-Type": []string{"text/plain"}},
			false,
		},
		{
			"Required \"driver\" property is missing",
			missingProp,
			jsonContentTypeHeader,
			false,
		},
		{
			"Not allowed \"type\" property",
			notAllowedByEnum,
			jsonContentTypeHeader,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest(
				http.MethodGet,
				"/driver_journeys?departureLat=0&departureLng=0&arrivalLat=0&arrivalLng=0&departureDate=1666014179&timeDelta=900&departureRadius=1&arrivalRadius=1",
				strings.NewReader(""),
			)
			panicIf(err)
			response := mockResponse(http.StatusOK, tc.body, tc.header)
			assertion := assertDriverJourneysFormat{request, response}
			assertionError := runSingleAssertion(t, assertion)
			if (assertionError == nil) != tc.expectNilError {
				t.Errorf("Wrong format response body should not be validated: %s",
					assertionError)
			}
		})
	}
}

func TestAssertAPICallSuccess(t *testing.T) {

	testCases := []struct {
		name           string
		apiCallError   error
		expectNilError bool
	}{
		{"nil error", nil, true},
		{"non nil error", errors.New(""), false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.apiCallError
			assertion := assertAPICallSuccess{err}
			assertionError := runSingleAssertion(t, assertion)
			if (assertionError == nil) != tc.expectNilError {
				t.Error("API call error is not handled as expected")
			}
		})
	}
}

func shouldHaveSingleAssertionResult(t *testing.T, ar []AssertionResult) {
	t.Helper()
	if len(ar) != 1 {
		t.Error("Each assertion should return only one AssertionResult")
	}
}

func TestDefaultAssertionAccu_Run(t *testing.T) {
	testCases := []struct {
		name                string
		assertions          []Assertion
		expectedNAssertions int
	}{
		{
			"Two success",
			[]Assertion{NoOpAssertion{}, NoOpAssertion{}},
			2,
		},
		{
			"Critic + success is not fatal",
			[]Assertion{Critic(NoOpAssertion{}), NoOpAssertion{}},
			2,
		},
		{
			"Critic + failure is fatal",
			[]Assertion{
				Critic(NoOpAssertion{errors.New("")}),
				NoOpAssertion{},
			},
			1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewAssertionAccu()
			a.Run(tc.assertions...)
			if len(a.storedAssertionResults) != tc.expectedNAssertions {
				t.Logf(
					"Got %d assertion executions, expected %d",
					len(a.storedAssertionResults),
					tc.expectedNAssertions,
				)
				t.Error("CriticAssertions are not handled as expected")
			}
		})
	}
}

// runSingleAssertion is a testing helper, which runs an assertion, and returns its underlying error (can
// be nil)
func runSingleAssertion(
	t *testing.T,
	assertion Assertion,
) error {
	t.Helper()
	a := NewAssertionAccu()
	a.Run(assertion)

	shouldHaveSingleAssertionResult(t, a.GetAssertionResults())

	return a.storedAssertionResults[0].err
}

func TestAssertRadius(t *testing.T) {
	var (
		coordsRef   = coords{46.1604531, -1.2219607} // reference
		coords900m  = coords{46.1613673, -1.2227555} // at ~900m from reference
		coords1100m = coords{46.1613679, -1.2086563} // at ~1100m from reference
	)

	testCases := []struct {
		name           string
		coordsResponse []coords
		radius         float32
		expectError    bool
	}{
		{
			name:           "no response",
			coordsResponse: []coords{},
			radius:         1,
			expectError:    false,
		},
		{
			name:           "1 inside radius 1km",
			coordsResponse: []coords{coords900m},
			radius:         1,
			expectError:    false,
		},
		{
			name:           "1 inside, 1 outside radius 1km",
			coordsResponse: []coords{coords900m, coords1100m},
			radius:         1,
			expectError:    true,
		},
		{
			name:           "2 inside, radius 1,2km",
			coordsResponse: []coords{coords900m, coords1100m},
			radius:         1.2,
			expectError:    false,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			params := client.GetDriverJourneysParams{
				DepartureRadius: &tc.radius,
				DepartureLat:    float32(coordsRef.lat),
				DepartureLng:    float32(coordsRef.lon),
			}
			request, err := client.NewGetDriverJourneysRequest("localhost:1323", &params)
			panicIf(err)

			responseObj := []client.DriverJourney{}
			for _, c := range tc.coordsResponse {
				responseObj = append(
					responseObj,
					client.DriverJourney{DriverDepartureLat: &c.lat, DriverDepartureLng: &c.lon},
				)
			}

			responseJSON, err := json.Marshal(responseObj)
			panicIf(err)
			response := mockResponse(200, string(responseJSON), nil)

			a := NewAssertionAccu()
			a.Run(assertDriverJourneysRadius{request, response, departure})

			results := a.GetAssertionResults()

			anyError := results[0].Unwrap() != nil
			if len(results) < 1 || anyError != tc.expectError {
				t.Error("Wrong behavior when asserting *radius query parameters")
			}
		})
	}
}
