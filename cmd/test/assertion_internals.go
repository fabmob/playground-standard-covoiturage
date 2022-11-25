package test

import "github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"

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

	// A string that summarizes the assertion
	assertionDescription string
}

// NewAssertionResult initializes an AssertionResult
func NewAssertionResult(err error, summary string) AssertionResult {
	return AssertionResult{
		err,
		summary,
	}
}

// Unwrap returns AssertionResult underlying error (possibly nil)
func (ar AssertionResult) Unwrap() error {
	return ar.err
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
	endpoint               endpoint.Info
}

// NewAssertionAccu inits a *DefaultAssertionAccu
func NewAssertionAccu() *DefaultAssertionAccu {
	return &DefaultAssertionAccu{
		storedAssertionResults: []AssertionResult{},
		endpoint:               endpoint.Info{},
	}
}

// Queue implements AssertionAccumulator.Queue.
func (a *DefaultAssertionAccu) Queue(assertions ...Assertion) {
	a.queuedAssertions = append(a.queuedAssertions, assertions...)
}

// ExecuteAll implements AssertionAccumulator.Run
func (a *DefaultAssertionAccu) ExecuteAll() {
	for _, assertion := range a.queuedAssertions {
		err := assertion.Execute()

		a.storedAssertionResults = append(
			a.storedAssertionResults,
			NewAssertionResult(err, assertion.Describe()),
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
