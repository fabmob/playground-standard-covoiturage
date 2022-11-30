package service

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/service/db"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	testassert "github.com/fabmob/playground-standard-covoiturage/cmd/test/assert"
	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const localServer = "http://localhost:1323"

func setupTestServer(
	db *db.Mock,
	request *http.Request,
) (*StdCovServerImpl, echo.Context, *httptest.ResponseRecorder) {

	var (
		e       = echo.New()
		rec     = httptest.NewRecorder()
		ctx     = e.NewContext(request, rec)
		handler = NewServerWithDB(db)
	)

	return handler, ctx, rec
}

func makeNTrips(n int) []api.Trip {
	trips := make([]api.Trip, 0, n)

	for i := 0; i < n; i++ {
		trips = append(trips, api.NewTrip())
	}

	return trips
}

func makeNDriverJourneys(n int) []api.DriverJourney {
	driverJourneys := make([]api.DriverJourney, 0, n)

	for i := 0; i < n; i++ {
		driverJourneys = append(driverJourneys, api.NewDriverJourney())
	}

	return driverJourneys
}

func makeNPassengerJourneys(n int) []api.PassengerJourney {
	passengerJourneys := make([]api.PassengerJourney, 0, n)

	for i := 0; i < n; i++ {
		passengerJourneys = append(passengerJourneys, api.NewPassengerJourney())
	}

	return passengerJourneys
}

func makeTripAtCoords(coordPickup, coordDrop util.Coord) api.Trip {
	t := api.NewTrip()
	updateTripCoords(&t, coordPickup, coordDrop)

	return t
}

func makeDriverJourneyAtCoords(coordPickup, coordDrop util.Coord) api.DriverJourney {
	dj := api.NewDriverJourney()
	updateTripCoords(&dj.Trip, coordPickup, coordDrop)

	return dj
}

func makePassengerJourneyAtCoords(coordPickup, coordDrop util.Coord) api.PassengerJourney {
	pj := api.NewPassengerJourney()
	updateTripCoords(&pj.Trip, coordPickup, coordDrop)

	return pj
}

func updateTripCoords(t *api.Trip, coordPickup, coordDrop util.Coord) {
	t.PassengerPickupLat = coordPickup.Lat
	t.PassengerPickupLng = coordPickup.Lon
	t.PassengerDropLat = coordDrop.Lat
	t.PassengerDropLng = coordDrop.Lon
}

func makeDriverJourneyAtDate(date int64) api.DriverJourney {
	dj := api.NewDriverJourney()
	dj.PassengerPickupDate = date

	return dj
}

func makePassengerJourneyAtDate(date int64) api.PassengerJourney {
	pj := api.NewPassengerJourney()
	pj.PassengerPickupDate = date

	return pj
}

func castDriverToPassenger(p *api.GetDriverJourneysParams) *api.GetPassengerJourneysParams {
	if p == nil {
		return nil
	}

	castedP := api.GetPassengerJourneysParams(*p)
	return &castedP
}

func makeParamsWithDepartureRadius(departureCoord util.Coord, departureRadius float32, driverOrPassenger string) api.GetJourneysParams {
	params := api.NewGetDriverJourneysParams(departureCoord, util.CoordIgnore, 0)
	params.DepartureRadius = &departureRadius

	if driverOrPassenger == "passenger" {
		return castDriverToPassenger(params)
	}

	return params
}

func makeParamsWithArrivalRadius(arrivalCoord util.Coord, arrivalRadius float32, driverOrPassenger string) api.GetJourneysParams {
	params := api.NewGetDriverJourneysParams(util.CoordIgnore, arrivalCoord, 0)
	params.ArrivalRadius = &arrivalRadius

	if driverOrPassenger == "passenger" {
		return castDriverToPassenger(params)
	}

	return params
}

func makeParamsWithTimeDelta(date int, driverOrPassenger string) api.GetJourneysParams {
	params := &api.GetDriverJourneysParams{}
	params.TimeDelta = &date

	if driverOrPassenger == "passenger" {
		return castDriverToPassenger(params)
	}

	return params
}

func makeRegularTripParamsWithTimeDelta(date int, driverOrPassenger string) *api.GetDriverRegularTripsParams {
	params := &api.GetDriverRegularTripsParams{}
	params.TimeDelta = &date

	return params
}

func makeParamsWithCount(count int, driverOrPassenger string) api.GetJourneysParams {
	params := &api.GetDriverJourneysParams{}
	params.Count = &count

	if driverOrPassenger == "passenger" {
		return castDriverToPassenger(params)
	}

	return params
}

// makeBooking sets status to "WAITING_CONFIRMATION" by default
func makeBooking(bookingID uuid.UUID) *api.Booking {
	return makeBookingWithStatus(bookingID, api.BookingStatusWAITINGCONFIRMATION)
}

func makeBookingWithStatus(bookingID uuid.UUID, status api.BookingStatus) *api.Booking {
	return &api.Booking{Id: bookingID, Status: status}
}

func makeCarpoolBookingEvent(eventID, bookingID uuid.UUID) *api.CarpoolBookingEvent {
	return makeCarpoolBookingEventWithStatus(eventID, bookingID,
		api.BookingStatusWAITINGCONFIRMATION)
}

func makeCarpoolBookingEventWithStatus(eventID, bookingID uuid.UUID, status api.BookingStatus) *api.CarpoolBookingEvent {
	booking := makeBookingWithStatus(bookingID, status)

	carpoolBookingEventData := api.CarpoolBookingEvent_Data{}
	err := carpoolBookingEventData.FromDriverCarpoolBooking(*booking.ToDriverCarpoolBooking())
	util.PanicIf(err)

	carpoolBookingEvent := &api.CarpoolBookingEvent{
		Data:    carpoolBookingEventData,
		Id:      eventID,
		IdToken: "",
	}

	return carpoolBookingEvent
}

func makeMessage(from api.User, to api.User) api.PostMessagesJSONBody {
	return api.PostMessagesJSONBody{
		From:                   from,
		To:                     to,
		Message:                "some message",
		RecipientCarpoolerType: "DRIVER",
	}
}

func makeUser(id, alias string) api.User {
	defaultOperator := "default.operator.com"
	return makeUserWithOperator(id, alias, defaultOperator)
}

func makeUserWithOperator(id, alias, operator string) api.User {
	return api.User{
		Id:       id,
		Alias:    alias,
		Operator: operator,
	}
}

// repUUID creates a reproducible UUID
func repUUID(seed int64) uuid.UUID {
	// generate random bytes
	rand.Seed(seed)
	randBytes := make([]byte, 16)
	rand.Read(randBytes)

	uuid, err := uuid.FromBytes(randBytes)
	util.PanicIf(err)

	return uuid
}

// NewBookingsByID populates a BookingsByID with given bookings. It does not
// test if booking is already set.
func NewBookingsByID(bookings ...*api.Booking) db.BookingsByID {
	var bookingsByID = db.BookingsByID{}

	for _, booking := range bookings {
		if booking == nil {
			panic(errors.New("attempt to insert nil booking"))
		}

		bookingsByID[booking.Id] = booking
	}

	return bookingsByID
}

func checkAssertionResults(t *testing.T, assertionResults []testassert.Result) {
	t.Helper()

	assert.Greater(t, len(assertionResults), 0)

	for _, ar := range assertionResults {
		if err := ar.Unwrap(); err != nil {
			t.Error(err)
		}
	}
}

func requestAll(t *testing.T, driverOrPassenger string) api.GetJourneysParams {
	t.Helper()

	var (
		largeTimeDelta = int(1e9)
		largeRadius    = float32(1e6)
	)

	switch driverOrPassenger {
	case "driver":
		params := api.GetDriverJourneysParams{}
		params.DepartureDate = 1e9
		params.TimeDelta = &largeTimeDelta
		params.DepartureRadius = &largeRadius
		params.ArrivalRadius = &largeRadius

		return &params

	case "passenger":
		params := api.GetPassengerJourneysParams{}
		params.DepartureDate = 1e9
		params.TimeDelta = &largeTimeDelta
		params.DepartureRadius = &largeRadius
		params.ArrivalRadius = &largeRadius

		return &params

	default:
		panic("invalid driverOrPassenger parameter")
	}
}

//////////////////////////////////////////////////////////

type apiTestHelper interface {
	makeRequest() (*http.Request, error)
	callAPI(*StdCovServerImpl, echo.Context) error
	testResponse(*http.Request, *http.Response, test.Flags) []testassert.Result
}

func testAPI(t *testing.T, a apiTestHelper, mockDB *db.Mock, flags test.Flags) {
	t.Helper()

	appendDataIfGenerated(t, mockDB)

	request, err := a.makeRequest()
	util.PanicIf(err)

	// Store server and endpoint information in request context
	server, endpointInfo, err := endpoint.FromRequest(request)
	util.PanicIf(err)

	requestCtx := endpoint.NewContext(context.Background(), server, endpointInfo)
	request = request.WithContext(requestCtx)

	// Setup testing server with response recorder
	handler, ctx, rec := setupTestServer(mockDB, request)

	// Read body without consuming it
	var body []byte
	if request.Body != nil {
		request.Body, err = test.ReusableReadCloser(ctx.Request().Body)
		util.PanicIf(err)
		body, err = io.ReadAll(request.Body)
		util.PanicIf(err)
	}

	// Make API Call
	err = a.callAPI(handler, ctx)
	util.PanicIf(err)

	appendCmdIfGenerated(t, request, flags, body)

	response := rec.Result()

	assertionResults := a.testResponse(request, response, flags)
	checkAssertionResults(t, assertionResults)
}

//////////////////////////////////////////////////////////

type postBookingsTestHelper struct {
	booking api.Booking
}

func (h postBookingsTestHelper) makeRequest() (*http.Request, error) {
	return api.NewPostBookingsRequest(localServer, h.booking)
}

func (h postBookingsTestHelper) callAPI(handler *StdCovServerImpl, ctx echo.Context) error {
	return handler.PostBookings(ctx)
}

func (h postBookingsTestHelper) testResponse(request *http.Request, response *http.Response, flags test.Flags) []testassert.Result {
	return test.TestPostBookingsResponse(request, response, flags)
}

func TestPostBookingsHelper(
	t *testing.T,
	mockDB *db.Mock,
	booking api.Booking,
	flags test.Flags,
) {
	testAPI(t, postBookingsTestHelper{booking}, mockDB, flags)
}

//////////////////////////////////////////////////////////

type postMessagesTestHelper struct {
	message api.PostMessagesJSONBody
}

func (h postMessagesTestHelper) makeRequest() (*http.Request, error) {
	return api.NewPostMessagesRequest(localServer, api.PostMessagesJSONRequestBody(h.message))
}

func (h postMessagesTestHelper) callAPI(handler *StdCovServerImpl, ctx echo.Context) error {
	return handler.PostMessages(ctx)
}

func (h postMessagesTestHelper) testResponse(request *http.Request, response *http.Response, flags test.Flags) []testassert.Result {
	return test.TestPostMessagesResponse(request, response, flags)
}

func TestPostMessagesHelper(
	t *testing.T,
	mockDB *db.Mock,
	message api.PostMessagesJSONBody,
	flags test.Flags,
) {
	testAPI(t, postMessagesTestHelper{message}, mockDB, flags)
}

//////////////////////////////////////////////////////////

type postBookingEventsTestHelper struct {
	bookingEvent api.CarpoolBookingEvent
}

func (h postBookingEventsTestHelper) makeRequest() (*http.Request, error) {
	return api.NewPostBookingEventsRequest(localServer, h.bookingEvent)
}

func (h postBookingEventsTestHelper) callAPI(handler *StdCovServerImpl, ctx echo.Context) error {
	return handler.PostBookingEvents(ctx)
}

func (h postBookingEventsTestHelper) testResponse(request *http.Request, response *http.Response, flags test.Flags) []testassert.Result {
	return test.TestPostBookingEventsResponse(request, response, flags)
}

func TestPostBookingEventsHelper(
	t *testing.T,
	mockDB *db.Mock,
	bookingEvent api.CarpoolBookingEvent,
	flags test.Flags,
) {
	testAPI(t, postBookingEventsTestHelper{bookingEvent}, mockDB, flags)
}

//////////////////////////////////////////////////////////

type getBookingsTestHelper struct {
	bookingID api.BookingId
}

func (h getBookingsTestHelper) makeRequest() (*http.Request, error) {
	return api.NewGetBookingsRequest(localServer, h.bookingID)
}

func (h getBookingsTestHelper) callAPI(handler *StdCovServerImpl, ctx echo.Context) error {
	return handler.GetBookings(ctx, h.bookingID)
}

func (h getBookingsTestHelper) testResponse(request *http.Request, response *http.Response, flags test.Flags) []testassert.Result {
	return test.TestGetBookingsResponse(request, response, flags)
}

func TestGetBookingsHelper(
	t *testing.T,
	mockDB *db.Mock,
	bookingID api.BookingId,
	flags test.Flags,
) {
	testAPI(t, getBookingsTestHelper{bookingID}, mockDB, flags)
}

//////////////////////////////////////////////////////////

type patchBookingsTestHelper struct {
	bookingID api.BookingId
	status    api.BookingStatus
}

func (h patchBookingsTestHelper) makeRequest() (*http.Request, error) {
	params := &api.PatchBookingsParams{Status: h.status}
	return api.NewPatchBookingsRequest(localServer, h.bookingID, params)
}

func (h patchBookingsTestHelper) callAPI(handler *StdCovServerImpl, ctx echo.Context) error {
	params := api.PatchBookingsParams{Status: h.status}
	return handler.PatchBookings(ctx, h.bookingID, params)
}

func (h patchBookingsTestHelper) testResponse(request *http.Request, response *http.Response, flags test.Flags) []testassert.Result {
	return test.TestPatchBookingsResponse(request, response, flags)
}

func TestPatchBookingsHelper(
	t *testing.T,
	mockDB *db.Mock,
	bookingID api.BookingId,
	status api.BookingStatus,
	flags test.Flags,
) {
	testAPI(t, patchBookingsTestHelper{bookingID, status}, mockDB, flags)
}

//////////////////////////////////////////////////////////

type getDriverJourneysTestHelper struct {
	params *api.GetDriverJourneysParams
}

func (h getDriverJourneysTestHelper) makeRequest() (*http.Request, error) {
	return api.NewGetDriverJourneysRequest(localServer, h.params)
}

func (h getDriverJourneysTestHelper) callAPI(handler *StdCovServerImpl, ctx echo.Context) error {
	return handler.GetDriverJourneys(ctx, *h.params)
}

func (h getDriverJourneysTestHelper) testResponse(request *http.Request, response *http.Response, flags test.Flags) []testassert.Result {
	return test.TestGetDriverJourneysResponse(request, response, flags)
}

func TestGetDriverJourneysHelper(
	t *testing.T,
	mockDB *db.Mock,
	params *api.GetDriverJourneysParams,
	flags test.Flags,
) {
	testAPI(t, getDriverJourneysTestHelper{params}, mockDB, flags)
}

//////////////////////////////////////////////////////////

type getPassengerJourneysTestHelper struct {
	params *api.GetPassengerJourneysParams
}

func (h getPassengerJourneysTestHelper) makeRequest() (*http.Request, error) {
	return api.NewGetPassengerJourneysRequest(localServer, h.params)
}

func (h getPassengerJourneysTestHelper) callAPI(handler *StdCovServerImpl, ctx echo.Context) error {
	return handler.GetPassengerJourneys(ctx, *h.params)
}

func (h getPassengerJourneysTestHelper) testResponse(request *http.Request, response *http.Response, flags test.Flags) []testassert.Result {
	return test.TestGetPassengerJourneysResponse(request, response, flags)
}

func TestGetPassengerJourneysHelper(
	t *testing.T,
	mockDB *db.Mock,
	params *api.GetPassengerJourneysParams,
	flags test.Flags,
) {
	testAPI(t, getPassengerJourneysTestHelper{params}, mockDB, flags)
}

//////////////////////////////////////////////////////////

type getDriverRegularTripsTestHelper struct {
	params *api.GetDriverRegularTripsParams
}

func (h getDriverRegularTripsTestHelper) makeRequest() (*http.Request, error) {
	return api.NewGetDriverRegularTripsRequest(localServer, h.params)
}

func (h getDriverRegularTripsTestHelper) callAPI(handler *StdCovServerImpl, ctx echo.Context) error {
	return handler.GetDriverRegularTrips(ctx, *h.params)
}

func (h getDriverRegularTripsTestHelper) testResponse(request *http.Request, response *http.Response, flags test.Flags) []testassert.Result {
	return test.TestGetDriverRegularTripsResponse(request, response, flags)
}

func TestGetDriverRegularTripsHelper(
	t *testing.T,
	mockDB *db.Mock,
	params *api.GetDriverRegularTripsParams,
	flags test.Flags,
) {
	testAPI(t, getDriverRegularTripsTestHelper{params}, mockDB, flags)
}

//////////////////////////////////////////////////////////

type getPassengerRegularTripsTestHelper struct {
	params *api.GetPassengerRegularTripsParams
}

func (h getPassengerRegularTripsTestHelper) makeRequest() (*http.Request, error) {
	return api.NewGetPassengerRegularTripsRequest(localServer, h.params)
}

func (h getPassengerRegularTripsTestHelper) callAPI(handler *StdCovServerImpl, ctx echo.Context) error {
	return handler.GetPassengerRegularTrips(ctx, *h.params)
}

func (h getPassengerRegularTripsTestHelper) testResponse(request *http.Request, response *http.Response, flags test.Flags) []testassert.Result {
	return test.TestGetPassengerRegularTripsResponse(request, response, flags)
}

func TestGetPassengerRegularTripsHelper(
	t *testing.T,
	mockDB *db.Mock,
	params *api.GetPassengerRegularTripsParams,
	flags test.Flags,
) {
	testAPI(t, getPassengerRegularTripsTestHelper{params}, mockDB, flags)
}

type driverJourneysTestCase struct {
	name                 string
	testParams           api.GetJourneysParams
	testData             []api.DriverJourney
	expectNonEmptyResult bool
}

type tripTestCase struct {
	name                 string
	testParams           api.JourneyOrTripPartialParams
	testData             []api.Trip
	expectNonEmptyResult bool
}

func promoteToDriverJourneysTestCase(t *testing.T, tc tripTestCase) driverJourneysTestCase {
	t.Helper()

	timeDelta := tc.testParams.GetTimeDelta()
	departureRadius := float32(tc.testParams.GetDepartureRadius())
	arrivalRadius := float32(tc.testParams.GetArrivalRadius())

	promotedParams := &api.GetDriverJourneysParams{
		DepartureLat:    float32(tc.testParams.GetDepartureLat()),
		DepartureLng:    float32(tc.testParams.GetDepartureLng()),
		ArrivalLat:      float32(tc.testParams.GetArrivalLat()),
		ArrivalLng:      float32(tc.testParams.GetArrivalLng()),
		TimeDelta:       &timeDelta,
		DepartureRadius: &departureRadius,
		ArrivalRadius:   &arrivalRadius,
		Count:           tc.testParams.GetCount(),
		DepartureDate:   0,
	}

	promotedData := []api.DriverJourney{}

	for _, d := range tc.testData {
		promotedTrip := api.NewDriverJourney()
		promotedTrip.Trip = d

		promotedData = append(promotedData, promotedTrip)
	}

	return driverJourneysTestCase{
		name:                 tc.name,
		testParams:           promotedParams,
		testData:             promotedData,
		expectNonEmptyResult: tc.expectNonEmptyResult,
	}
}
