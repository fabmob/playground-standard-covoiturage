package test

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gitlab.com/multi/stdcov-api-test/cmd/api"
	"gitlab.com/multi/stdcov-api-test/cmd/util"
	"gitlab.com/multi/stdcov-api-test/cmd/validate"
)

// An Assertion is a unit test that can be executed and that can describe
// itself
type Assertion interface {
	Execute() error
	Describe() string
}

// An AssertionResult stores data and metadata about the result of a single assertion
type AssertionResult struct {
	// Error, if any
	err error

	// Endpoint under test
	endpoint Endpoint

	// A string that summarizes the assertion
	assertionDescription string
}

// NewAssertionResult initializes an AssertionResult
func NewAssertionResult(err error, endpointPath, endpointMethod, summary string) AssertionResult {
	return AssertionResult{
		err,
		Endpoint{endpointMethod, endpointPath},
		summary,
	}
}

// Unwrap returns AssertionResult underlying error (possibly nil)
func (ar AssertionResult) Unwrap() error {
	return ar.err
}

// String implements Stringer interface.
// Formats the AssertionResult nicely in one line (no linebreak).
func (ar AssertionResult) String() string {

	err := ar.Unwrap()

	var symbol string
	if err != nil {
		symbol = "ERROR ❌"
	} else {
		symbol = "OK ✅"
	}

	resStr := fmt.Sprintf(
		"%7s %-35s  %-35s",
		symbol,
		ar.endpoint,
		ar.assertionDescription,
	)
	if err != nil {
		resStr += fmt.Sprintf("\n%5s %s", "", err)
	}
	return resStr
}

/////////////////////////////////////////////////////////////

// A CriticAssertion is an Assertion, which success is critic for the
// execution of subsequent assertions.
type CriticAssertion struct {
	Assertion
}

// Critic converts an Assertion into a CriticAssertion
func Critic(a Assertion) CriticAssertion {
	return CriticAssertion{a}
}

// An AssertionAccumulator can run assertions, store and retrieve the
// corresponding AssertionResults
type AssertionAccumulator interface {
	// Queue adds assertion to the queue for later execution
	Queue(...Assertion)

	// ExecuteAll executes assertions in sequence and stores the results.
	// If a CriticAssertion fails, execution is interrupted.
	ExecuteAll()

	// GetAssertionResults returns all results of executed assertions
	GetAssertionResults() []AssertionResult
}

/////////////////////////////////////////////////////////////

// DefaultAssertionAccu implements Asserter interface
type DefaultAssertionAccu struct {
	queuedAssertions       []Assertion
	storedAssertionResults []AssertionResult
	endpoint               Endpoint
}

// NewAssertionAccu inits a *DefaultAssertionAccu
func NewAssertionAccu() *DefaultAssertionAccu {
	return &DefaultAssertionAccu{
		storedAssertionResults: []AssertionResult{},
		endpoint:               Endpoint{},
	}
}

func (a *DefaultAssertionAccu) Queue(assertions ...Assertion) {
	for _, assertion := range assertions {
		a.queuedAssertions = append(a.queuedAssertions, assertion)
	}
}

// ExecuteAll implements AssertionAccumulator.Run
func (a *DefaultAssertionAccu) ExecuteAll() {
	for _, assertion := range a.queuedAssertions {
		err := assertion.Execute()

		a.storedAssertionResults = append(
			a.storedAssertionResults,
			NewAssertionResult(err, a.endpoint.Path, a.endpoint.Method,
				assertion.Describe()),
		)
		_, critic := assertion.(CriticAssertion)
		fatal := (critic && err != nil)
		if fatal {
			return
		}
	}
}

// GetAssertionResults implements AssertionAccumulator.GetAssertionResults
func (a *DefaultAssertionAccu) GetAssertionResults() []AssertionResult {
	return a.storedAssertionResults
}

/////////////////////////////////////////////////////////////
// Helper functions
/////////////////////////////////////////////////////////////

// AssertAPICallSuccess checks if requesting an endpoint returned an error
func AssertAPICallSuccess(a AssertionAccumulator, err error) {
	assertion := assertAPICallSuccess{err}
	a.Queue(assertion)
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
	assertion := assertDriverJourneysFormat{request, response}
	a.Queue(assertion)
}

func CriticAssertDriverJourneysFormat(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := Critic(assertDriverJourneysFormat{request, response})
	a.Queue(assertion)
}

func AssertDriverJourneysDepartureRadius(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := assertDriverJourneysRadius{request, response, departure}
	a.Queue(assertion)
}

func AssertDriverJourneysArrivalRadius(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := assertDriverJourneysRadius{request, response, arrival}
	a.Queue(assertion)
}

func AssertDriverJourneysNotEmpty(a AssertionAccumulator, response *http.Response) {
	assertion := assertDriverJourneysNotEmpty{response}
	a.Queue(assertion)
}

func CriticAssertDriverJourneysNotEmpty(a AssertionAccumulator, response *http.Response) {
	assertion := Critic(assertDriverJourneysNotEmpty{response})
	a.Queue(assertion)
}

func AssertDriverJourneysTimeDelta(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := assertDriverJourneysTimeDelta{request, response}
	a.Queue(assertion)
}

func AssertDriverJourneysCount(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := assertDriverJourneysCount{request, response}
	a.Queue(assertion)
}

func AssertUniqueIDs(a AssertionAccumulator, response *http.Response) {
	assertion := assertUniqueIDs{response}
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

type assertDriverJourneysFormat struct {
	request  *http.Request
	response *http.Response
}

func (a assertDriverJourneysFormat) Execute() error {
	err := validate.Response(a.request, a.response)
	return err
}

func (a assertDriverJourneysFormat) Describe() string {
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
	queryParams, err := api.ParseGetDriverJourneysRequest(a.request)
	if err != nil {
		return failedParsing("request", err)
	}
	coordsQuery := getQueryCoord(a.departureOrArrival, queryParams)
	// As different distance computations may give different distances, we apply
	// a safety margin
	radius := getQueryRadiusOrDefault(a.departureOrArrival, queryParams)
	safetyMarginPercent := 1.
	radiusWithMargin := radius * (1. + safetyMarginPercent/100)

	// Parse response
	driverJourneys, err := api.ParseGetDriverJourneysOKResponse(a.response)
	if err != nil {
		return failedParsing("response", err)
	}

	for _, dj := range driverJourneys {
		coordsResponse, err := getResponseCoord(a.departureOrArrival, dj)
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

type assertDriverJourneysNotEmpty struct {
	response *http.Response
}

func (a assertDriverJourneysNotEmpty) Execute() error {
	driverJourneys, err := api.ParseGetDriverJourneysOKResponse(a.response)
	if err != nil {
		return failedParsing("response", err)
	}
	if len(driverJourneys) == 0 {
		return errors.New("empty response not accepted with \"disallowEmpty\" option")
	}
	return nil
}

func (a assertDriverJourneysNotEmpty) Describe() string {
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
