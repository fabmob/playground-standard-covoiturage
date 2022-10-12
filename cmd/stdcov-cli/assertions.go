package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-cli/client"
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
		Endpoint{endpointPath, endpointMethod},
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

// An AssertionAccumulator can run assertions, store and retrieve the
// corresponding AssertionResults
type AssertionAccumulator interface {
	// Run executes the given Assertion and stores the result
	Run(Assertion)

	// GetAssertionResults returns all results of executed assertions
	GetAssertionResults() []AssertionResult
	// LastAssertionHasError returns whether the last assertion returned an
	// error
	LastAssertionHasError() bool
}

/////////////////////////////////////////////////////////////

// DefaultAssertionAccu implements Asserter interface
type DefaultAssertionAccu struct {
	storedAssertionResults []AssertionResult
	endpoint               Endpoint
}

// NewDefaultAsserter inits a *DefaultAsserter
func NewDefaultAsserter() *DefaultAssertionAccu {
	return &DefaultAssertionAccu{
		storedAssertionResults: []AssertionResult{},
		endpoint:               Endpoint{},
	}
}

// Run implements Asserter.Run
func (a *DefaultAssertionAccu) Run(assertion Assertion) {
	err := assertion.Execute()
	a.storedAssertionResults = append(
		a.storedAssertionResults,
		NewAssertionResult(err, a.endpoint.path, a.endpoint.method,
			assertion.Describe()),
	)
}

// GetAssertionResults implements Asserter.GetAssertionResults
func (a *DefaultAssertionAccu) GetAssertionResults() []AssertionResult {
	return a.storedAssertionResults
}

// LastAssertionHasError implements Asserter.LastAssertionHasError
func (a *DefaultAssertionAccu) LastAssertionHasError() bool {
	ar := a.storedAssertionResults
	if len(ar) == 0 {
		panic("Trying to access inexistant or empty []AssertionError")
	}
	return ar[len(ar)-1].Unwrap() != nil
}

/////////////////////////////////////////////////////////////
// Helper functions
/////////////////////////////////////////////////////////////

// AssertAPICallSuccess checks if requesting an endpoint returned an error
func AssertAPICallSuccess(a AssertionAccumulator, err error) {
	assertion := assertAPICallSuccess{err}
	a.Run(assertion)
}

// AssertStatusCode checks if a given response has an expected status code
/* AssertStatusCode(*http.Response, int) */
func AssertStatusCode(a AssertionAccumulator, resp *http.Response, statusCode int) {
	assertion := assertStatusCode{resp, statusCode}
	a.Run(assertion)
}

// AssertStatusCodeOK checks if a given response has status 200 OK
func AssertStatusCodeOK(a AssertionAccumulator, resp *http.Response) {
	AssertStatusCode(a, resp, http.StatusOK)
}

// AssertHeaderContains checks if a given key is present in header, with associated
// value
func AssertHeaderContains(a AssertionAccumulator, resp *http.Response, key, value string) {
	assertion := assertHeaderContains{resp, key, value}
	a.Run(assertion)
}

// AssertDriverJourneysFormat checks if the response data of
// /driver_journeys call has the expected format
func AssertDriverJourneysFormat(a AssertionAccumulator, response *http.Response) {
	assertion := assertDriverJourneysFormat{response}
	a.Run(assertion)
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
	response *http.Response
}

func (a assertDriverJourneysFormat) Execute() error {
	a.response.Header["Content-Type"] = []string{"json"}
	_, err := client.ParseGetDriverJourneysResponse(a.response)
	return err
}

func (a assertDriverJourneysFormat) Describe() string {
	return "assertDriverJourneysFormat"
}