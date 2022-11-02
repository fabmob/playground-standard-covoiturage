package test

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	tld "github.com/jpillora/go-tld"
	"github.com/pkg/errors"
	"gitlab.com/multi/stdcov-api-test/cmd/util"
	"gitlab.com/multi/stdcov-api-test/cmd/validate"
)

/////////////////////////////////////////////////////////////
// Helper functions
/////////////////////////////////////////////////////////////

// CheckAPICallSuccess checks if requesting an endpoint returned an error
func CheckAPICallSuccess(err error) AssertionResult {
	assertion := assertAPICallSuccess{err}
	return NewAssertionResult(assertion.Execute(), assertion.Describe())
}

// AssertStatusCode checks if a given response has an expected status code
/* AssertStatusCode(*http.Response, int) */
func AssertStatusCode(a AssertionAccumulator, resp *http.Response, statusCode int) {
	assertion := assertStatusCode{resp, statusCode}
	a.Queue(assertion)
}

// AssertStatusCodeOK checks if a given response has status 200 OK
func AssertStatusCodeOK(a AssertionAccumulator, resp *http.Response) {
	AssertStatusCode(a, resp, http.StatusOK)
}

// AssertHeaderContains checks if a given key is present in header, with associated
// value
func AssertHeaderContains(a AssertionAccumulator, resp *http.Response, key, value string) {
	assertion := assertHeaderContains{resp, key, value}
	a.Queue(assertion)
}

// AssertFormat checks if the response data of
// /driver_journeys call has the expected format
func AssertFormat(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := assertFormat{request, response}
	a.Queue(assertion)
}

// CriticAssertFormat checks if the response data of
// /driver_journeys call has the expected format. A failure prevents the
// following assertions to be executed.
func CriticAssertFormat(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := Critic(assertFormat{request, response})
	a.Queue(assertion)
}

// AssertJourneysDepartureRadius checks that the response data respect
// the "departureRadius" query parameter
func AssertJourneysDepartureRadius(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := assertJourneysRadius{request, response, departure}
	a.Queue(assertion)
}

// AssertJourneysArrivalRadius checks that the response data respect
// the "arrivalRadius" query parameter
func AssertJourneysArrivalRadius(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := assertJourneysRadius{request, response, arrival}
	a.Queue(assertion)
}

// AssertArrayNotEmpty checks that the response is not empty
func AssertArrayNotEmpty(a AssertionAccumulator, response *http.Response) {
	assertion := assertArrayNotEmpty{response}
	a.Queue(assertion)
}

// CriticAssertArrayNotEmpty checks that the response is not empty. A
// failure prevents the following assertions to be executed.
func CriticAssertArrayNotEmpty(a AssertionAccumulator, response *http.Response) {
	assertion := Critic(assertArrayNotEmpty{response})
	a.Queue(assertion)
}

// AssertJourneysTimeDelta checks that the response data respect the
// "timeDelta" query parameter
func AssertJourneysTimeDelta(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := assertJourneysTimeDelta{request, response}
	a.Queue(assertion)
}

// AssertJourneysCount checks that the response data respect the "count"
// query parameter
func AssertJourneysCount(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := assertJourneysCount{request, response}
	a.Queue(assertion)
}

// AssertUniqueIDs checks that all driverJourneys IDs, if they exist, are
// unique.
func AssertUniqueIDs(a AssertionAccumulator, response *http.Response) {
	assertion := assertUniqueIDs{response}
	a.Queue(assertion)
}

// AssertOperatorFieldFormat checks that the response data has well formated
// "operator" field
func AssertOperatorFieldFormat(a AssertionAccumulator, response *http.Response) {
	assertion := assertOperatorFieldFormat{response}
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
	return "assertAPICallSuccess"
}

/////////////////////////////////////////////////////////////

type assertStatusCode struct {
	resp       *http.Response
	statusCode int
}

func (a assertStatusCode) Execute() error {
	expected := a.statusCode
	got := a.resp.StatusCode
	if expected != got {
		return (errors.Errorf("Expected status code %d, got %d", expected, got))
	}
	return nil
}

func (a assertStatusCode) Describe() string {
	return "assertStatusCode " + strconv.Itoa(a.statusCode)
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
	return "assertheader " + a.key + ":" + a.value
}

/////////////////////////////////////////////////////////////

type assertFormat struct {
	request  *http.Request
	response *http.Response
}

func (a assertFormat) Execute() error {
	err := validate.Response(a.request, a.response)
	return err
}

func (a assertFormat) Describe() string {
	return "assertDriverJourneysFormat"
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
	responseObjects, err := parseArrayResponse(a.response)
	if err != nil {
		return failedParsing("response", err)
	}

	for _, obj := range responseObjects {
		coordsResponse, err := getResponseCoord(a.departureOrArrival, obj)
		if err != nil {
			return err
		}
		dist := util.Distance(coordsResponse, coordsQuery)
		if dist > radius {
			return fmt.Errorf("a driver journey does not comply to maximum '%s' distance to query parameters", a.departureOrArrival)
		}
	}
	return nil
}

func (a assertJourneysRadius) Describe() string {
	return fmt.Sprintf("assert %s", a.departureOrArrival)
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
		return errors.New("empty response not accepted with \"disallowEmpty\" option")
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

	responseObjects, err := parseArrayResponse(a.response)
	if err != nil {
		return failedParsing("response", err)
	}

	for _, obj := range responseObjects {
		pickupDate, err := getResponsePickupDate(obj)
		if err != nil {
			return failedParsing("response", err)
		}
		if math.Abs(float64(pickupDate)-float64(departureDate)) >
			float64(timeDelta) {
			return errors.New("a driver journey does not comply to timeDelta query parameter")
		}
	}
	return nil
}

func (a assertJourneysTimeDelta) Describe() string {

	return "assert timeDelta"
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
	objs, err := parseArrayResponse(a.response)
	if err != nil {
		return err
	}
	if count != -1 {
		if len(objs) > count {
			return errors.New("the number of returned journeys exceeds the query count parameter")
		}
	}
	return nil
}

func (a assertJourneysCount) Describe() string {
	return "assert count"
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
	return "assert operator field format"
}
