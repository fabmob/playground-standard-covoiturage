package assert

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var localServer = "http://localhost:1323"

// emptyRequest returns an empty *http.Request to the endpoint
func emptyRequest(e endpoint.Info) *http.Request {
	request, err := http.NewRequest(e.Method, localServer+e.Path, nil)
	util.PanicIf(err)

	request, err = endpoint.AddEndpointContext(request)
	util.PanicIf(err)

	return request
}

// A NopAssertion returns stored error when executed
type NopAssertion struct{ error }

// Execute implements Assertion interface
func (n NopAssertion) Execute() error {
	return n.error
}

// Describe implements Assertion interface
func (NopAssertion) Describe() string {
	return "No assertion"
}

// makeJourneyRequestWithRadius is a test helper that creates a request,
// either for GET /driver_journey or GET /passenger_journey, with given
// coords, either for departure or arrival, and given radius.
func makeJourneyRequestWithRadius(
	t *testing.T,
	coord util.Coord,
	radius float32,
	departureOrArrival departureOrArrival,
	driverOrPassenger string, // "driver" or "passenger"
) *http.Request {

	t.Helper()

	var params api.GetDriverJourneysParams

	switch departureOrArrival {
	case departure:

		params = api.GetDriverJourneysParams{
			DepartureRadius: &radius,
			DepartureLat:    float32(coord.Lat),
			DepartureLng:    float32(coord.Lon),
		}

	case arrival:

		params = api.GetDriverJourneysParams{
			ArrivalRadius: &radius,
			ArrivalLat:    float32(coord.Lat),
			ArrivalLng:    float32(coord.Lon),
		}

	default:
		panic(errors.New("wrong value in test: departureOrArrival"))
	}

	var request *http.Request
	var err error

	switch driverOrPassenger {
	case "driver":
		castedParams := api.GetDriverJourneysParams(params)
		request, err = api.NewGetDriverJourneysRequest(localServer, &castedParams)
		util.PanicIf(err)
	case "passenger":
		castedParams := api.GetPassengerJourneysParams(params)
		request, err = api.NewGetPassengerJourneysRequest(localServer, &castedParams)
		util.PanicIf(err)
	case "default":
		panic(errors.New("wrong value in test: driverOrPassenger"))
	}

	return request
}

func makeJourneysResponse(t *testing.T, coords []util.Coord, departureOrArrival departureOrArrival, driverOrPassenger string) *http.Response {
	t.Helper()

	var (
		response                *http.Response
		driverJourneyObjects    = []api.DriverJourney{}
		passengerJourneyObjects = []api.PassengerJourney{}
	)

	for _, c := range coords {
		var trip api.Trip

		if departureOrArrival == departure {
			trip.PassengerPickupLat = c.Lat
			trip.PassengerPickupLng = c.Lon
		} else {
			trip.PassengerDropLat = c.Lat
			trip.PassengerDropLng = c.Lon
		}

		switch driverOrPassenger {
		case "driver":
			driverJourneyObjects = append(
				driverJourneyObjects,
				api.DriverJourney{DriverTrip: api.DriverTrip{Trip: trip}},
			)

		case "passenger":
			passengerJourneyObjects = append(
				passengerJourneyObjects,
				api.PassengerJourney{PassengerTrip: api.PassengerTrip{Trip: trip}},
			)

		default:
			panic(errors.New("wrong value in test: driverOrPassenger"))
		}
	}

	switch driverOrPassenger {
	case "driver":
		response = mockBodyResponse(interface{}(driverJourneyObjects))

	case "passenger":
		response = mockBodyResponse(interface{}(passengerJourneyObjects))

	default:
		panic(errors.New("wrong value in test: driverOrPassenger"))
	}

	return response
}

// errAsExpected returns if the error is as expected.
// (expectError = false <=> err == nil)
func errAsExpected(err error, expectError bool) bool {
	hasError := (err != nil)
	return hasError == expectError
}

// mockResponse returns a mock response with given statusCode, body, and
// header. If headers are `nil` default headers with "Content-Type: json" are
// used.
func mockResponse(
	statusCode int,
	body string,
	header http.Header,
) *http.Response {

	if header == nil {
		header = make(http.Header)
		header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	return &http.Response{
		Status:        http.StatusText(statusCode),
		StatusCode:    statusCode,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
		Header:        header,
	}
}

func mockStatusResponse(statusCode int) *http.Response {
	return mockResponse(statusCode, "", nil)
}

func MockOKStatusResponse() *http.Response {
	return mockStatusResponse(http.StatusOK)
}

func mockBodyResponse(responseObj interface{}) *http.Response {
	responseJSON, err := json.Marshal(responseObj)
	util.PanicIf(err)

	return mockResponse(200, string(responseJSON), nil)
}

func ShouldHaveSingleAssertionResult(t *testing.T, ar []Result) {
	t.Helper()
	if len(ar) != 1 {
		t.Error("Each assertion should return only one AssertionResult")
	}
}
