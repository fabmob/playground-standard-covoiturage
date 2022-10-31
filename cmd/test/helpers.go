package test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"gitlab.com/multi/stdcov-api-test/cmd/api"
	"gitlab.com/multi/stdcov-api-test/cmd/util"
)

func getQueryRadius(departureOrArrival departureOrArrival, req *http.Request) float64 {
	const DefaultRadius float64 = 1

	radiusStr := req.URL.Query().Get(string(departureOrArrival))

	var radius float64
	if radiusStr == "" {
		return DefaultRadius
	}
	radius, err := strconv.ParseFloat(radiusStr, 64)
	panicIf(err) // Should never happen it request format is validated
	return radius
}

// getQueryCoord extracts departure or arrival coordinates from
// queryParameters
func getQueryCoord(departureOrArrival departureOrArrival, request *http.Request) util.Coord {
	var latParam, lonParam string
	switch departureOrArrival {
	case departure:
		latParam = "departureLat"
		lonParam = "departureLng"
	case arrival:
		latParam = "arrivalLat"
		lonParam = "arrivalLng"
	}
	latStr := request.URL.Query().Get(latParam)
	lat, _ := strconv.ParseFloat(latStr, 64)
	lonStr := request.URL.Query().Get(lonParam)
	lon, _ := strconv.ParseFloat(lonStr, 64)
	coordQuery := util.Coord{Lat: lat, Lon: lon}
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
