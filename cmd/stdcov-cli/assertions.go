package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
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
	// Run executes assertions in sequence and stores the results.
	// If a CriticAssertion fails, execution is interrupted.
	Run(...Assertion)

	// GetAssertionResults returns all results of executed assertions
	GetAssertionResults() []AssertionResult
}

/////////////////////////////////////////////////////////////

// DefaultAssertionAccu implements Asserter interface
type DefaultAssertionAccu struct {
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

// Run implements AssertionAccumulator.Run
func (a *DefaultAssertionAccu) Run(assertions ...Assertion) {
	for _, assertion := range assertions {
		err := assertion.Execute()

		a.storedAssertionResults = append(
			a.storedAssertionResults,
			NewAssertionResult(err, a.endpoint.path, a.endpoint.method,
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
func AssertDriverJourneysFormat(a AssertionAccumulator, request *http.Request, response *http.Response) {
	assertion := assertDriverJourneysFormat{request, response}
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
