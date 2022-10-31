package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"gitlab.com/multi/stdcov-api-test/cmd/util"
)

/////////////////////////////////////////////////////////////
// Query parameter extraction
/////////////////////////////////////////////////////////////

func getQueryRadius(departureOrArrival departureOrArrival, req *http.Request) (float64, error) {
	const DefaultRadius float64 = 1
	return parseQueryFloatParamWithDefault(
		req,
		string(departureOrArrival),
		DefaultRadius,
	)
}

// getQueryCoord extracts departure or arrival coordinates from
// queryParameters
func getQueryCoord(departureOrArrival departureOrArrival, request *http.Request) (util.Coord, error) {
	var latParam, lonParam string
	switch departureOrArrival {
	case departure:
		latParam = "departureLat"
		lonParam = "departureLng"
	case arrival:
		latParam = "arrivalLat"
		lonParam = "arrivalLng"
	}
	lat, err := parseQueryFloatParam(request, latParam)
	if err != nil {
		return util.Coord{}, err
	}
	lon, err := parseQueryFloatParam(request, lonParam)
	if err != nil {
		return util.Coord{}, err
	}
	coordQuery := util.Coord{Lat: lat, Lon: lon}
	return coordQuery, nil
}

func parseQueryFloatParam(request *http.Request, paramName string) (float64, error) {
	paramStr := request.URL.Query().Get(paramName)
	param, err := strconv.ParseFloat(paramStr, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"%s could not be properly parsed as float in query (%w)",
			paramStr,
			err,
		)
	}
	return param, nil
}

func parseQueryFloatParamWithDefault(request *http.Request, paramName string, defaultValue float64) (float64, error) {
	paramStr := request.URL.Query().Get(paramName)
	if paramStr == "" {
		return defaultValue, nil
	}
	param, err := strconv.ParseFloat(paramStr, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"%s could not be properly parsed as float in query (%w)",
			paramStr,
			err,
		)
	}
	return param, nil
}

/////////////////////////////////////////////////////////////
// Response body extraction
/////////////////////////////////////////////////////////////

// getResponseCoord extracts departure or arrival coordinates from
// a json.RawMessage, e.g. as returned by parseArrayResponse. Fails if response has no such coordinates.
func getResponseCoord(departureOrArrival departureOrArrival, obj json.RawMessage) (util.Coord, error) {
	var coordResponse util.Coord

	switch departureOrArrival {
	case departure:
		type PassengerPickupCoord struct {
			PassengerPickupLat float64 `json:"passengerPickupLat"`
			PassengerPickupLng float64 `json:"passengerPickupLng"`
		}
		var passengerPickupCoord PassengerPickupCoord
		err := json.Unmarshal(obj, &passengerPickupCoord)
		if err != nil {
			return util.Coord{}, err
		}

		coordResponse = util.Coord{
			Lat: passengerPickupCoord.PassengerPickupLat,
			Lon: passengerPickupCoord.PassengerPickupLng,
		}
	case arrival:
		type PassengerDropCoord struct {
			PassengerDropLat float64 `json:"passengerDropLat"`
			PassengerDropLng float64 `json:"passengerDropLng"`
		}
		var passengerDropCoord PassengerDropCoord
		err := json.Unmarshal(obj, &passengerDropCoord)
		if err != nil {
			return util.Coord{}, err
		}

		coordResponse = util.Coord{
			Lat: passengerDropCoord.PassengerDropLat,
			Lon: passengerDropCoord.PassengerDropLng,
		}
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

/////////////////////////////////////////////////////////////

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
