// Package assert defines atomic assertions to be used in tests
package assert

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
	tld "github.com/jpillora/go-tld"
	"github.com/pkg/errors"
)

/////////////////////////////////////////////////////////////
// Exported Assertion functions
/////////////////////////////////////////////////////////////

// CheckAPICallSuccess checks if requesting an endpoint returned an error
func CheckAPICallSuccess(err error) Result {
	assertion := assertAPICallSuccess{err}
	return NewAssertionResult(assertion.Execute(), assertion.Describe())
}

// StatusCode checks if a given response has an expected status code
/* StatusCode(*http.Response, int) */
func StatusCode(a Accumulator, resp *http.Response, statusCode int) {
	assertion := assertStatusCode{resp, statusCode}
	a.Queue(assertion)
}

// StatusCodeOK checks if a given response has status 200 OK
func StatusCodeOK(a Accumulator, resp *http.Response) {
	StatusCode(a, resp, http.StatusOK)
}

// HeaderContains checks if a given key is present in header, with associated
// value
func HeaderContains(a Accumulator, resp *http.Response, key, value string) {
	assertion := assertHeaderContains{resp, key, value}
	a.Queue(assertion)
}

// Format checks if the response data has the expected format
func Format(a Accumulator, request *http.Request, response *http.Response) {
	assertion := assertFormat{request, response}
	a.Queue(assertion)
}

// CriticFormat is the same as Format, but a failure prevents the
// following assertions to be executed.
func CriticFormat(a Accumulator, request *http.Request, response *http.Response) {
	assertion := Critic(assertFormat{request, response})
	a.Queue(assertion)
}

// JourneysDepartureRadius checks that the response data respect
// the "departureRadius" query parameter
func JourneysDepartureRadius(a Accumulator, request *http.Request, response *http.Response) {
	assertion := assertJourneysRadius{request, response, departure}
	a.Queue(assertion)
}

// JourneysArrivalRadius checks that the response data respect
// the "arrivalRadius" query parameter
func JourneysArrivalRadius(a Accumulator, request *http.Request, response *http.Response) {
	assertion := assertJourneysRadius{request, response, arrival}
	a.Queue(assertion)
}

// ArrayNotEmpty checks that the response is not empty
func ArrayNotEmpty(a Accumulator, response *http.Response) {
	assertion := assertArrayNotEmpty{response}
	a.Queue(assertion)
}

// CriticArrayNotEmpty checks that the response is not empty. A
// failure prevents the following assertions to be executed.
func CriticArrayNotEmpty(a Accumulator, response *http.Response) {
	assertion := Critic(assertArrayNotEmpty{response})
	a.Queue(assertion)
}

// JourneysTimeDelta checks that the response data respect the
// "timeDelta" query parameter
func JourneysTimeDelta(a Accumulator, request *http.Request, response *http.Response) {
	assertion := assertJourneysTimeDelta{request, response}
	a.Queue(assertion)
}

// JourneysCount checks that the response data respect the "count"
// query parameter
func JourneysCount(a Accumulator, request *http.Request, response *http.Response) {
	assertion := assertJourneysCount{request, response}
	a.Queue(assertion)
}

// UniqueIDs checks that all IDs (property "id"), if they exist, are
// unique.
func UniqueIDs(a Accumulator, response *http.Response) {
	assertion := assertUniqueIDs{response}
	a.Queue(assertion)
}

// OperatorFieldFormat checks that the response data has well formated
// "operator" field
func OperatorFieldFormat(a Accumulator, response *http.Response) {
	assertion := assertOperatorFieldFormat{response}
	a.Queue(assertion)
}

func BookingStatus(a Accumulator, response *http.Response, expectedStatus string) {
	assertion := assertBookingStatus{response, expectedStatus}
	a.Queue(assertion)
}

/////////////////////////////////////////////////////////////

type assertAPICallSuccess struct {
	apiCallErr error
}

func (a assertAPICallSuccess) Execute() error {
	if a.apiCallErr != nil {
		return a.apiCallErr
	}

	return nil
}

func (a assertAPICallSuccess) Describe() string {
	return "assert API call success"
}

/////////////////////////////////////////////////////////////

type assertStatusCode struct {
	resp       *http.Response
	statusCode int
}

func (a assertStatusCode) Execute() error {
	expected, got := a.statusCode, a.resp.StatusCode
	if expected != got {
		return (errors.Errorf("Expected status code %d, got %d", expected, got))
	}

	return nil
}

func (a assertStatusCode) Describe() string {
	return fmt.Sprintf("assert status code %d", a.statusCode)
}

/////////////////////////////////////////////////////////////

type assertHeaderContains struct {
	resp       *http.Response
	key, value string
}

func (a assertHeaderContains) Execute() error {
	if val, ok := a.resp.Header[a.key]; !ok {
		return errors.Errorf("expected header %s, which is missing", a.key)
	} else if len(val[0]) < 1 || !strings.Contains(val[0], a.value) {
		return errors.Errorf(
			"expected value %s for header %s, got %s",
			a.value,
			a.key,
			val,
		)
	} else {
		return nil
	}
}

func (a assertHeaderContains) Describe() string {
	return fmt.Sprintf("assert header %s:%s", a.key, a.value)
}

/////////////////////////////////////////////////////////////

type assertFormat struct {
	request  *http.Request
	response *http.Response
}

func (a assertFormat) Execute() error {
	err := validateResponse(a.request, a.response)
	return err
}

func (a assertFormat) Describe() string {
	return "assert format"
}

/////////////////////////////////////////////////////////////

type departureOrArrival string

const (
	departure departureOrArrival = "departureRadius"
	arrival   departureOrArrival = "arrivalRadius"
)

// assertJourneysRadius expects that response format has been validated
type assertJourneysRadius struct {
	request            *http.Request
	response           *http.Response
	departureOrArrival departureOrArrival
}

func (a assertJourneysRadius) Execute() error {
	// Parse request
	coordsQuery, err := getQueryCoord(a.departureOrArrival, a.request)
	if err != nil {
		return failedParsing("request", err)
	}

	radius, err := getQueryRadius(a.departureOrArrival, a.request)
	if err != nil {
		return failedParsing("request", err)
	}

	// As different distance computations may give different distances, we apply
	// a safety margin
	safetyMarginPercent := 1.
	radius = radius * (1. + safetyMarginPercent/100)

	// Parse response
	objsWithRadius, err := parseArrayResponse(a.response)
	if err != nil {
		return failedParsing("response", err)
	}

	for _, objWithRadius := range objsWithRadius {
		coordsResponse, err := getResponseCoord(a.departureOrArrival, objWithRadius)
		if err != nil {
			return err
		}

		dist := util.Distance(coordsResponse, coordsQuery)
		if dist > radius {
			return fmt.Errorf("a journey does not comply to maximum '%s' distance to query parameters", a.departureOrArrival)
		}
	}

	return nil
}

func (a assertJourneysRadius) Describe() string {
	return fmt.Sprintf("assert query parameter \"%s\"", a.departureOrArrival)
}

/////////////////////////////////////////////////////////////

type assertArrayNotEmpty struct {
	response *http.Response
}

func (a assertArrayNotEmpty) Execute() error {
	array, err := parseArrayResponse(a.response)
	if err != nil {
		return failedParsing("response", err)
	}

	if len(array) == 0 {
		return errors.New("empty response not accepted with \"expectNonEmpty\" option")
	}

	return nil
}

func (a assertArrayNotEmpty) Describe() string {
	return "assert response not empty"
}

/////////////////////////////////////////////////////////////

type assertJourneysTimeDelta struct {
	request  *http.Request
	response *http.Response
}

func (a assertJourneysTimeDelta) Execute() error {
	timeDelta, err := getQueryTimeDelta(a.request)
	if err != nil {
		return failedParsing("request", err)
	}

	departureDate, err := getQueryDeparturDate(a.request)
	if err != nil {
		return failedParsing("request", err)
	}

	objsWithTimeDelta, err := parseArrayResponse(a.response)
	if err != nil {
		return failedParsing("response", err)
	}

	for _, objWithTimeDelta := range objsWithTimeDelta {
		pickupDate, err := getResponsePickupDate(objWithTimeDelta)
		if err != nil {
			return failedParsing("response", err)
		}

		if math.Abs(float64(pickupDate)-float64(departureDate)) >
			float64(timeDelta) {
			return errors.New("a journey does not comply to timeDelta query parameter")
		}
	}

	return nil
}

func (a assertJourneysTimeDelta) Describe() string {
	return "assert query parameter \"timeDelta\""
}

/////////////////////////////////////////////////////////////

type assertJourneysCount struct {
	request  *http.Request
	response *http.Response
}

func (a assertJourneysCount) Execute() error {
	count, err := getQueryCount(a.request)
	if err != nil {
		return failedParsing("request", err)
	}

	objsWithCount, err := parseArrayResponse(a.response)
	if err != nil {
		return err
	}

	if count != -1 {
		if len(objsWithCount) > count {
			return errors.New("the number of returned journeys exceeds the query count parameter")
		}
	}

	return nil
}

func (a assertJourneysCount) Describe() string {
	return "assert query parameter \"count\""
}

/////////////////////////////////////////////////////////////

type assertUniqueIDs struct {
	response *http.Response
}

func (a assertUniqueIDs) Execute() error {
	objsWithID, err := parseArrayResponse(a.response)
	if err != nil {
		return failedParsing("response", err)
	}

	ids := map[string]bool{}

	for _, objWithID := range objsWithID {
		idPtr, err := getResponseID(objWithID)
		if err != nil {
			return failedParsing("response", err)
		}

		if idPtr != nil {
			id := *idPtr

			_, idDuplicate := ids[id]
			if idDuplicate {
				return errors.New("IDs should be unique")
			}

			ids[id] = true
		}
	}

	return nil
}

func (a assertUniqueIDs) Describe() string {
	return "assert unique ids"
}

/////////////////////////////////////////////////////////////

type assertOperatorFieldFormat struct {
	response *http.Response
}

func (a assertOperatorFieldFormat) Execute() error {
	objsWithOperator, err := parseArrayResponse(a.response)
	if err != nil {
		return err
	}

	for _, objWithOperator := range objsWithOperator {
		operator, err := getResponseOperator(objWithOperator)
		if err != nil {
			return failedParsing("response", err)
		}

		if err := validateOperator(operator); err != nil {
			return err
		}
	}

	return nil
}

func validateOperator(operator string) error {
	uri, err := tld.Parse("https://" + operator)
	if err != nil {
		return fmt.Errorf("wrong operator field format: %w", err)
	}

	if uri.Host == "" || uri.Path != "" || uri.User != nil || uri.RawQuery != "" {
		return fmt.Errorf("wrong operator field format")
	}

	return nil
}

func (a assertOperatorFieldFormat) Describe() string {
	return "assert response property \"operator\""
}

/////////////////////////////////////////////////////////////

type assertBookingStatus struct {
	response       *http.Response
	expectedStatus string
}

func (a assertBookingStatus) Execute() error {
	bodyBytes, err := io.ReadAll(a.response.Body)
	if err != nil {
		return err
	}

	status, err := getResponseStatus(json.RawMessage(bodyBytes))
	if err != nil {
		return err
	}

	if status != a.expectedStatus {
		return fmt.Errorf(
			"expected booking status %s, got %s", a.expectedStatus, status,
		)
	}

	return nil
}

func (a assertBookingStatus) Describe() string {
	return fmt.Sprintf("assert booking status %s", a.expectedStatus)
}
