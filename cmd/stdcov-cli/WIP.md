
- Test that the right assertions are called with fake `Asserter`
- Factorize assertions_test.go thanks to new `Assertion` interface



Vocabulary :

- API
- Endpoint
- Each endpoint undergoes several **tests**. The collection of tests is called 
  **test suite**.
- A test can build on minimal building blocks called `Assertion`s
- Each test returns a collection of `AssertionResult`s. 
- The organized set of all `AssertionResults` is called a `Report`


