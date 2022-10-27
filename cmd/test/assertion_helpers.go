package test

import (
	"bytes"
	"fmt"
	"io"

	"gitlab.com/multi/stdcov-api-test/cmd/api"
	"gitlab.com/multi/stdcov-api-test/cmd/util"
)

// getQueryCoord extracts departure or arrival coordinates from
// queryParameters
func getQueryCoord(departureOrArrival departureOrArrival, queryParams *api.GetDriverJourneysParams) util.Coord {
	var coordQuery util.Coord
	switch departureOrArrival {
	case departure:
		coordQuery = util.Coord{Lat: float64(queryParams.DepartureLat), Lon: float64(queryParams.DepartureLng)}
	case arrival:
		coordQuery = util.Coord{Lat: float64(queryParams.ArrivalLat), Lon: float64(queryParams.ArrivalLng)}
	}
	return coordQuery
}

// getResponseCoord extracts departure or arrival coordinates from
// driverJourney object. Fails if required coordinates are missing.
func getResponseCoord(departureOrArrival departureOrArrival, driverJourney api.DriverJourney) (util.Coord, error) {
	var coordResponse util.Coord
	switch departureOrArrival {
	case departure:
		coordResponse = util.Coord{Lat: driverJourney.PassengerPickupLat, Lon: driverJourney.PassengerPickupLng}
	case arrival:
		coordResponse = util.Coord{Lat: driverJourney.PassengerDropLat, Lon: driverJourney.PassengerDropLng}
	}
	return coordResponse, nil
}

// getQueryRadiusOrDefault returns departureRadius er arrivalRadius query parameter
// (depending on departureOrArrival), or the default value if missing
func getQueryRadiusOrDefault(departureOrArrival departureOrArrival, queryParams *api.GetDriverJourneysParams) float64 {
	const DefaultRadius float32 = 1
	var radiusPtr *float32
	switch departureOrArrival {
	case departure:
		radiusPtr = queryParams.DepartureRadius
	case arrival:
		radiusPtr = queryParams.ArrivalRadius
	}
	var radius float32
	if radiusPtr != nil {
		radius = *radiusPtr
	} else {
		radius = DefaultRadius
	}
	return float64(radius)
}

// failedParsing wraps a parsing error with additional details
func failedParsing(responseOrRequest string, err error) error {
	return fmt.Errorf(
		"internal error while parsing %s:%w",
		responseOrRequest,
		err,
	)
}

type reusableReadCloser struct {
	io.ReadCloser
	readBuf *bytes.Buffer
	backBuf *bytes.Buffer
}

// ReusableReadCloser wraps a io.ReadCloser so that it can be read and closed as
// many times as needed
func ReusableReadCloser(r io.ReadCloser) io.ReadCloser {
	readBuf := bytes.Buffer{}
	readBuf.ReadFrom(r) // error handling ignored for brevity
	backBuf := bytes.Buffer{}

	return reusableReadCloser{
		io.NopCloser(io.TeeReader(&readBuf, &backBuf)),
		&readBuf,
		&backBuf,
	}
}

func (r reusableReadCloser) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	if err == io.EOF {
		r.reset()
	}
	return n, err
}

func (r reusableReadCloser) reset() {
	io.Copy(r.readBuf, r.backBuf) // nolint: errcheck
}

func (r reusableReadCloser) Close() error {
	return nil
}
