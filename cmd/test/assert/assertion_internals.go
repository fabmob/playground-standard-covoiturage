package assert

import "github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"

// An Assertion is a unit test that can be executed and that can describe
// itself
type Assertion interface {
	Execute() error
	Describe() string
}

// An Result stores data and metadata about the result of a single assertion
type Result struct {
	// Error, if any
	Err error

	// A string that summarizes the assertion
	AssertionDescription string
}

// NewAssertionResult initializes an AssertionResult
func NewAssertionResult(err error, summary string) Result {
	return Result{
		err,
		summary,
	}
}

// Unwrap returns AssertionResult underlying error (possibly nil)
func (ar Result) Unwrap() error {
	return ar.Err
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

// An Accumulator can run assertions, store and retrieve the
// corresponding AssertionResults
type Accumulator interface {
	// Queue adds assertion to the queue for later execution
	Queue(...Assertion)

	// ExecuteAll executes assertions in sequence and stores the results.
	// If a CriticAssertion fails, execution is interrupted.
	ExecuteAll()

	// GetAssertionResults returns all results of executed assertions
	GetAssertionResults() []Result
}

/////////////////////////////////////////////////////////////

// DefaultAccumulator implements Asserter interface
type DefaultAccumulator struct {
	queuedAssertions       []Assertion
	storedAssertionResults []Result
	endpoint               endpoint.Info
}

// NewAccumulator inits a *DefaultAssertionAccu
func NewAccumulator() *DefaultAccumulator {
	return &DefaultAccumulator{
		storedAssertionResults: []Result{},
		endpoint:               endpoint.Info{},
	}
}

// Queue implements Accumulator.Queue.
func (a *DefaultAccumulator) Queue(assertions ...Assertion) {
	a.queuedAssertions = append(a.queuedAssertions, assertions...)
}

// ExecuteAll implements Accumulator.Run
func (a *DefaultAccumulator) ExecuteAll() {
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

// GetAssertionResults implements Accumulator.GetAssertionResults
func (a *DefaultAccumulator) GetAssertionResults() []Result {
	return a.storedAssertionResults
}
