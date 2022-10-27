package test

import (
	"bytes"
	"fmt"
	"io"

	"github.com/umahmood/haversine"
	"gitlab.com/multi/stdcov-api-test/cmd/api"
)

type coords struct {
	lat float64
	lon float64
}

func distanceKm(coords1, coords2 coords) float64 {
	c1 := haversine.Coord{Lat: coords1.lat, Lon: coords1.lon}
	c2 := haversine.Coord{Lat: coords2.lat, Lon: coords2.lon}
	_, dist := haversine.Distance(c1, c2)
	return dist
}

// getQueryCoords extracts departure or arrival coordinates from
// queryParameters
func getQueryCoords(departureOrArrival departureOrArrival, queryParams *api.GetDriverJourneysParams) coords {
	var coordsQuery coords
	switch departureOrArrival {
	case departure:
		coordsQuery = coords{float64(queryParams.DepartureLat), float64(queryParams.DepartureLng)}
	case arrival:
		coordsQuery = coords{float64(queryParams.ArrivalLat), float64(queryParams.ArrivalLng)}
	}
	return coordsQuery
}

// getResponseCoords extracts departure or arrival coordinates from
// driverJourney object. Fails if required coordinates are missing.
func getResponseCoords(departureOrArrival departureOrArrival, driverJourney api.DriverJourney) (coords, error) {
	var coordsResponse coords
	switch departureOrArrival {
	case departure:
		coordsResponse = coords{driverJourney.PassengerPickupLat, driverJourney.PassengerPickupLng}
	case arrival:
		coordsResponse = coords{driverJourney.PassengerDropLat, driverJourney.PassengerDropLng}
	}
	return coordsResponse, nil
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
