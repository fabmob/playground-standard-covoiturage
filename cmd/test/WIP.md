
- Test that the right assertions are called with fake `Asserter`
- Factorize assertions_test.go thanks to new `Assertion` interface
- Check that URL option is not empty (or set default to server)
- assertDriverJourneysFormat should not modify the response object in-place

Server:

- Validate request with OapiRequestValidator middleware

Possible assertions driver journeys:

- "the carpooling operator SHOULD return in priority the most relevant 
  results. The measure of relevance is left to the discretion of the 
  carpooling operator."
- unique ids, same operator fields, operator fields format
- weburl required if deeplink supported.
- long-lat in France ?

Possible assertions booking object:

- 404 if missing, 200 otherwise
- PAYING = amount required
- Booking by API = PAYING required
- driverJourneyID, passengerJourneyId (how to check ? "If the booking is made 
  after a search, the MaaS platform SHOULD recall the journey IDs.")
- Unique (no way to test) and UUID format ID 
- Driver and passenger operator formats

Ideas:

- Each unit test data set is on a separate day. It is therefore possible to 
  consolidate unit test data into a single default dataset on which to run the 
  test queries.
- A flag "--failIfEmpty" which does not accept empty responses

Validation:

- Validate response data from file on import ! 

Issues:
 
- price "type" should be required
- no error code field while message says "Error code can be among".
- id description for CarpoolBookingEvent is wrong
- Tags should separate MaaS from operator endpoints in openapi Spec
- Reuse CarpoolBooking object in Booking object
- DriverCarpoolBooking / FormatCarpoolBooking id format uuid
