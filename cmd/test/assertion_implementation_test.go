package test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
	"github.com/labstack/echo/v4"
)

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

		assertionError := singleAssertionError(t, assertion)
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
		echo.HeaderContentType: {echo.MIMEApplicationJSON},
	}
	headerContentTypeJSONWithCharset := http.Header{
		echo.HeaderContentType: {echo.MIMEApplicationJSONCharsetUTF8},
	}
	headerContentTypeForm := http.Header{
		echo.HeaderContentType: {echo.MIMEMultipartForm},
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
			echo.HeaderContentType,
			echo.MIMEApplicationJSON,
			false,
		},
		{
			"json Content-Type header",
			headerContentTypeJSON,
			echo.HeaderContentType,
			echo.MIMEApplicationJSON,
			true,
		},
		{
			"json Content-Type header with charset",
			headerContentTypeJSONWithCharset,
			echo.HeaderContentType,
			echo.MIMEApplicationJSON,
			true,
		},
		{
			"json Content-Type header",
			headerContentTypeJSON,
			"Server",
			echo.MIMEApplicationJSON,
			false,
		},
		{
			"wrong Content-Type header",
			headerContentTypeForm,
			echo.HeaderContentType,
			echo.MIMEApplicationJSON,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := mockResponse(http.StatusOK, "", tc.header)
			assertion := assertHeaderContains{r, tc.testKey, tc.testValue}
			assertionError := singleAssertionError(t, assertion)
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

	marshalBody := func(dj interface{}) string {
		bodyBytes, _ := json.Marshal(dj)
		return string(bodyBytes)
	}

	invalidDJ := api.NewDriverJourney()
	invalidDJ.Type = "Not allowed"

	var (
		// Test requests
		driverJourneysRequest    = emptyRequest(GetDriverJourneysEndpoint)
		passengerJourneysRequest = emptyRequest(GetPassengerJourneysEndpoint)

		// Test bodies
		emptyDriverJourneysBody    = marshalBody([]api.DriverJourney{})
		singleDriverJourneyBody    = marshalBody([]api.DriverJourney{api.NewDriverJourney()})
		singlePassengerJourneyBody = marshalBody([]api.PassengerJourney{api.NewPassengerJourney()})
		notAllowedByEnum           = marshalBody([]api.DriverJourney{invalidDJ})
		missingProp                = `[
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
	)

	jsonContentTypeHeader := http.Header{echo.HeaderContentType: []string{echo.MIMEApplicationJSON}}

	testCases := []struct {
		name           string
		request        *http.Request
		body           string
		header         http.Header
		expectNilError bool
	}{
		{
			"Not JSON",
			driverJourneysRequest,
			"Hello, world!",
			jsonContentTypeHeader,
			false,
		},
		{
			"Empty []DriverJourney JSON",
			driverJourneysRequest,
			emptyDriverJourneysBody,
			jsonContentTypeHeader,
			true,
		},
		{
			"Non-empty []DriverJourney JSON",
			driverJourneysRequest,
			singleDriverJourneyBody,
			jsonContentTypeHeader,
			true,
		},
		{
			"Non-empty []PassengerJourney JSON",
			passengerJourneysRequest,
			singlePassengerJourneyBody,
			jsonContentTypeHeader,
			true,
		},
		{
			"Other content type",
			driverJourneysRequest,
			"Hello, world!",
			http.Header{echo.HeaderContentType: []string{echo.MIMETextPlain}},
			false,
		},
		{
			"Required \"driver\" property is missing",
			driverJourneysRequest,
			missingProp,
			jsonContentTypeHeader,
			false,
		},
		{
			"Not allowed \"type\" property",
			driverJourneysRequest,
			notAllowedByEnum,
			jsonContentTypeHeader,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := tc.request
			response := mockResponse(http.StatusOK, tc.body, tc.header)
			assertion := assertFormat{request, response}
			assertionError := singleAssertionError(t, assertion)
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
			assertionError := singleAssertionError(t, assertion)
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
			[]Assertion{NopAssertion{}, NopAssertion{}},
			2,
		},
		{
			"Critic + success is not fatal",
			[]Assertion{Critic(NopAssertion{}), NopAssertion{}},
			2,
		},
		{
			"Critic + failure is fatal",
			[]Assertion{
				Critic(NopAssertion{errors.New("")}),
				NopAssertion{},
			},
			1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewAssertionAccu()
			a.Queue(tc.assertions...)
			a.ExecuteAll()
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

// singleAssertionError is a testing helper, which runs an assertion, and returns its underlying error (can
// be nil)
func singleAssertionError(
	t *testing.T,
	assertion Assertion,
) error {
	t.Helper()

	var a = NewAssertionAccu()

	a.Queue(assertion)
	a.ExecuteAll()

	shouldHaveSingleAssertionResult(t, a.GetAssertionResults())

	return a.storedAssertionResults[0].err
}

func TestAssertRadius(t *testing.T) {
	var (
		coordsRef   = util.Coord{Lat: 46.1604531, Lon: -1.2219607} // reference
		coords900m  = util.Coord{Lat: 46.1613442, Lon: -1.2103736} // at ~900m from reference
		coords1100m = util.Coord{Lat: 46.1613679, Lon: -1.2086563} // at ~1100m from reference
	)

	testCases := []struct {
		name               string
		departureOrArrival departureOrArrival
		driverOrPassenger  string
		coordRequest       util.Coord
		coordsResponse     []util.Coord
		radius             float32
		expectError        bool
	}{
		{
			name:               "no response",
			departureOrArrival: departure,
			driverOrPassenger:  "driver",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{},
			radius:             1,
			expectError:        false,
		},
		{
			name:               "1 inside radius 1km",
			departureOrArrival: departure,
			driverOrPassenger:  "driver",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{coords900m},
			radius:             1,
			expectError:        false,
		},
		{
			name:               "1 inside, 1 outside radius 1km",
			departureOrArrival: departure,
			driverOrPassenger:  "driver",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{coords900m, coords1100m},
			radius:             1,
			expectError:        true,
		},
		{
			name:               "2 inside, radius 1,2km",
			departureOrArrival: departure,
			driverOrPassenger:  "driver",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{coords900m, coords1100m},
			radius:             1.2,
			expectError:        false,
		},
		{
			name:               "1 inside, other reference, radius 0.5km",
			departureOrArrival: departure,
			driverOrPassenger:  "driver",
			coordRequest:       coords900m,
			coordsResponse:     []util.Coord{coords1100m},
			radius:             0.5,
			expectError:        false,
		},
		{
			name:               "no response",
			departureOrArrival: arrival,
			driverOrPassenger:  "driver",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{},
			radius:             1,
			expectError:        false,
		},
		{
			name:               "1 inside radius 1km",
			departureOrArrival: arrival,
			driverOrPassenger:  "driver",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{coords900m},
			radius:             1,
			expectError:        false,
		},
		{
			name:               "1 inside, 1 outside radius 1km",
			departureOrArrival: arrival,
			driverOrPassenger:  "driver",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{coords900m, coords1100m},
			radius:             1,
			expectError:        true,
		},
		{
			name:               "2 inside, radius 1,2km",
			departureOrArrival: arrival,
			driverOrPassenger:  "driver",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{coords900m, coords1100m},
			radius:             1.2,
			expectError:        false,
		},
		{
			name:               "1 inside, other reference, radius 0.5km",
			departureOrArrival: arrival,
			driverOrPassenger:  "driver",
			coordRequest:       coords900m,
			coordsResponse:     []util.Coord{coords1100m},
			radius:             0.5,
			expectError:        false,
		},
		{
			name:               "no response",
			departureOrArrival: departure,
			driverOrPassenger:  "passenger",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{},
			radius:             1,
			expectError:        false,
		},
		{
			name:               "1 inside radius 1km",
			departureOrArrival: departure,
			driverOrPassenger:  "passenger",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{coords900m},
			radius:             1,
			expectError:        false,
		},
		{
			name:               "1 inside, 1 outside radius 1km",
			departureOrArrival: departure,
			driverOrPassenger:  "passenger",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{coords900m, coords1100m},
			radius:             1,
			expectError:        true,
		},
		{
			name:               "2 inside, radius 1,2km",
			departureOrArrival: departure,
			driverOrPassenger:  "passenger",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{coords900m, coords1100m},
			radius:             1.2,
			expectError:        false,
		},
		{
			name:               "1 inside, other reference, radius 0.5km",
			departureOrArrival: departure,
			driverOrPassenger:  "passenger",
			coordRequest:       coords900m,
			coordsResponse:     []util.Coord{coords1100m},
			radius:             0.5,
			expectError:        false,
		},
		{
			name:               "no response",
			departureOrArrival: arrival,
			driverOrPassenger:  "passenger",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{},
			radius:             1,
			expectError:        false,
		},
		{
			name:               "1 inside radius 1km",
			departureOrArrival: arrival,
			driverOrPassenger:  "passenger",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{coords900m},
			radius:             1,
			expectError:        false,
		},
		{
			name:               "1 inside, 1 outside radius 1km",
			departureOrArrival: arrival,
			driverOrPassenger:  "passenger",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{coords900m, coords1100m},
			radius:             1,
			expectError:        true,
		},
		{
			name:               "2 inside, radius 1,2km",
			departureOrArrival: arrival,
			driverOrPassenger:  "passenger",
			coordRequest:       coordsRef,
			coordsResponse:     []util.Coord{coords900m, coords1100m},
			radius:             1.2,
			expectError:        false,
		},
		{
			name:               "1 inside, other reference, radius 0.5km",
			departureOrArrival: arrival,
			driverOrPassenger:  "passenger",
			coordRequest:       coords900m,
			coordsResponse:     []util.Coord{coords1100m},
			radius:             0.5,
			expectError:        false,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			request := makeJourneyRequestWithRadius(
				t,
				tc.coordRequest,
				tc.radius,
				tc.departureOrArrival,
				tc.driverOrPassenger,
			)
			response := makeJourneysResponse(
				t,
				tc.coordsResponse,
				tc.departureOrArrival,
				tc.driverOrPassenger,
			)

			t.Log(response)
			err := singleAssertionError(
				t,
				assertJourneysRadius{request, response, tc.departureOrArrival},
			)

			if !errAsExpected(err, tc.expectError) {
				t.Log(err)
				t.Error("Wrong behavior when asserting *radius query parameters")
			}
		})
	}
}

func TestAssertNotEmpty(t *testing.T) {
	testCases := []struct {
		name         string
		responseData []interface{}
		expectError  bool
	}{
		{"empty whatever", []interface{}{}, true},
		{"non empty driver journeys", []interface{}{api.DriverJourney{}}, false},
		{"non empty passenger journeys", []interface{}{api.PassengerJourney{}}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := mockBodyResponse(tc.responseData)
			err := singleAssertionError(t, assertArrayNotEmpty{response})
			if !errAsExpected(err, tc.expectError) {
				t.Fail()
			}
		})
	}
}

func TestAssertUniqueIDs(t *testing.T) {
	var (
		id1          = "1"
		id1duplicate = "1"
		id2          = "2"
	)

	testCases := []struct {
		name        string
		ids         []*string
		expectError bool
	}{
		{"no id", []*string{nil, nil}, false},
		{"unique ids", []*string{&id1, &id2}, false},
		{"duplicate id", []*string{&id1, &id2, &id1duplicate}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			responseData := make([]api.DriverJourney, 0, len(tc.ids))

			for _, id := range tc.ids {
				dj := api.NewDriverJourney()
				dj.Id = id
				responseData = append(responseData, dj)
			}

			response := mockBodyResponse(responseData)

			err := singleAssertionError(t, assertUniqueIDs{response})
			if !errAsExpected(err, tc.expectError) {
				t.Fail()
			}
		})
	}
}

func TestValidateOperator(t *testing.T) {
	testCases := []struct {
		operator string
		valid    bool
	}{
		{"operator.com", true},
		{"operator.fr", true},
		{"carpooling.com", true},
		{"subdomain.operator.com", true},
		{"subdomain.subdomain.operator.co.uk", true},
		{"", false},
		{"random", false},
		{"https://operator.com", false},
		{"operator.com/", false},
		{"operator.com/path", false},
		{"/some/path", false},
	}

	for _, tc := range testCases {
		if err := validateOperator(tc.operator); !errAsExpected(err, !tc.valid) {
			t.Logf("Operator: %s, Expected to be valid: %t", tc.operator, tc.valid)
			t.Logf("Error: %s", err)
			t.Fail()
		}
	}
}

func TestExpectedBookingStatus(t *testing.T) {
	testCases := []struct {
		bookingStatus  api.BookingStatus
		expectedStatus api.BookingStatus
		expectError    bool
	}{
		{
			api.BookingStatusCANCELLED,
			api.BookingStatusCANCELLED,
			false,
		}, {
			api.BookingStatusCANCELLED,
			api.BookingStatusVALIDATED,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run("Test case", func(t *testing.T) {
			statusObj := struct{ Status string }{string(tc.bookingStatus)}
			response := mockBodyResponse(statusObj)

			err := singleAssertionError(t, assertBookingStatus{response, string(tc.expectedStatus)})
			if !errAsExpected(err, tc.expectError) {
				t.Logf("Expected status %s, got %s", tc.expectedStatus, tc.bookingStatus)
				t.Fail()
			}
		})
	}
}
