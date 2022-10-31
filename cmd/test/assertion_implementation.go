package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"

	tld "github.com/jpillora/go-tld"
	"github.com/pkg/errors"
	"gitlab.com/multi/stdcov-api-test/cmd/api"
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

// AssertDriverJourneysFormat checks if the response data of
// /driver_journeys call has the expected format
func AssertDriverJourneysFormat(a AssertionAccumulator, request *http.Request, response *http.Response) {
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

// AssertDriverJourneysDepartureRadius checks that the response data respect
// the "departureRadius" query parameter
func AssertDriverJourneysDepartureRadius(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := assertDriverJourneysRadius{request, response, departure}
	a.Queue(assertion)
}

// AssertDriverJourneysArrivalRadius checks that the response data respect
// the "arrivalRadius" query parameter
func AssertDriverJourneysArrivalRadius(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := assertDriverJourneysRadius{request, response, arrival}
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

// AssertDriverJourneysTimeDelta checks that the response data respect the
// "timeDelta" query parameter
func AssertDriverJourneysTimeDelta(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := assertDriverJourneysTimeDelta{request, response}
	a.Queue(assertion)
}

// AssertDriverJourneysCount checks that the response data respect the "count"
// query parameter
func AssertDriverJourneysCount(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := assertDriverJourneysCount{request, response}
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

// assertDriverJourneysRadius expects that response format has been validated
type assertDriverJourneysRadius struct {
	request            *http.Request
	response           *http.Response
	departureOrArrival departureOrArrival
}

func (a assertDriverJourneysRadius) Execute() error {
	// Parse request
	coordsQuery, err := getQueryCoord(a.departureOrArrival, a.request)
	if err != nil {
		return failedParsing("request", err)
	}
	// As different distance computations may give different distances, we apply
	// a safety margin
	radius, err := getQueryRadius(a.departureOrArrival, a.request)
	if err != nil {
		return failedParsing("request", err)
	}

	safetyMarginPercent := 1.
	radiusWithMargin := radius * (1. + safetyMarginPercent/100)

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
		if dist > radiusWithMargin {
			return fmt.Errorf("a driver journey does not comply to maximum '%s' distance to query parameters", a.departureOrArrival)
		}
	}
	return nil
}

func (a assertDriverJourneysRadius) Describe() string {
	return fmt.Sprintf("assert %s", a.departureOrArrival)
}

/////////////////////////////////////////////////////////////

type assertArrayNotEmpty struct {
	response *http.Response
}

// parseArrayOKResponse parses an array of any type, keeping array elements as
// json.RawMessage
func parseArrayResponse(rsp *http.Response) ([]json.RawMessage, error) {

	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}
	var dest []json.RawMessage
	if err := json.Unmarshal(bodyBytes, &dest); err != nil {
		return nil, err
	}
	return dest, nil
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

type assertDriverJourneysTimeDelta struct {
	request  *http.Request
	response *http.Response
}

func (a assertDriverJourneysTimeDelta) Execute() error {
	params, err := api.ParseGetDriverJourneysRequest(a.request)
	if err != nil {
		return err
	}
	driverJourneys, err := api.ParseGetDriverJourneysOKResponse(a.response)
	if err != nil {
		return err
	}
	for _, dj := range driverJourneys {
		if math.Abs(float64(dj.PassengerPickupDate)-float64(params.DepartureDate)) >
			float64(params.GetTimeDelta()) {
			return errors.New("a driver journey does not comply to timeDelta query parameter")
		}
	}
	return nil
}

func (a assertDriverJourneysTimeDelta) Describe() string {

	return "assert timeDelta"
}

/////////////////////////////////////////////////////////////

type assertDriverJourneysCount struct {
	request  *http.Request
	response *http.Response
}

func (a assertDriverJourneysCount) Execute() error {
	params, err := api.ParseGetDriverJourneysRequest(a.request)
	if err != nil {
		return err
	}
	driverJourneys, err := api.ParseGetDriverJourneysOKResponse(a.response)
	if err != nil {
		return err
	}
	if params.Count != nil {
		expectedMaxCount := params.Count
		if len(driverJourneys) > *expectedMaxCount {
			return errors.New("the number of returned driver journeys exceeds the query count parameter")
		}
	}
	return nil
}

func (a assertDriverJourneysCount) Describe() string {
	return "assert count"
}

/////////////////////////////////////////////////////////////

type assertUniqueIDs struct {
	response *http.Response
}

func (a assertUniqueIDs) Execute() error {
	driverJourneys, err := api.ParseGetDriverJourneysOKResponse(a.response)
	if err != nil {
		return err
	}
	ids := map[string]bool{}
	for _, dj := range driverJourneys {
		if dj.Id != nil {
			id := *dj.Id
			_, idDuplicate := ids[id]
			if idDuplicate {
				return errors.New("IDs should be unique")
			}
			ids[*dj.Id] = true
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
	driverJourneys, err := api.ParseGetDriverJourneysOKResponse(a.response)
	if err != nil {
		return err
	}
	for _, dj := range driverJourneys {
		if err := validateOperator(dj.Operator); err != nil {
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
	if uri.Path != "" || uri.User != nil || uri.RawQuery != "" {
		return fmt.Errorf("wrong operator field format")
	}
	return nil
}

func (a assertOperatorFieldFormat) Describe() string {
	return "assert operator field format"
}
