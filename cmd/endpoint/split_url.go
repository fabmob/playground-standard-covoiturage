package endpoint

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

var allEndpoints = []Info{
	GetDriverJourneys,
	GetPassengerJourneys,
	GetDriverRegularTrips,
	GetPassengerRegularTrips,
	PostBookingEvents,
	PostMessages,
	PostBookings,
	PatchBookings,
	GetBookings,
	GetStatus,
}

// splitURL tries to guess the server, and returns server and path in case of
// success.
func splitURL(method, URL string) (Server, Info, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", Info{}, err
	}

	removeQuery(u)

	for _, endpoint := range allEndpoints {

		suffix := knownEndpointSuffix(u.String(), endpoint)

		if endpoint.Method == method && suffix != "" {
			server := strings.TrimSuffix(u.String(), suffix)
			return Server(server), endpoint, nil
		}
	}

	return Server(""), Info{}, fmt.Errorf(
		"did not recognize the endpoint with method %s in %s",
		method,
		u,
	)
}

func removeQuery(u *url.URL) {
	u.RawQuery = ""
	u.Fragment = ""
}

// knownEndpointSuffix returns the complete endpoint suffix (including path
// parameter) if the endpoint path is recognized, an empty string
// otherwise.
func knownEndpointSuffix(url string, endpoint Info) string {
	var param string

	if endpoint.HasPathParam {
		url, param = path.Split(url)
		url = ensureNoTrailingSlash(url)
		param = ensureLeadingSlash(param)
	}

	if strings.HasSuffix(url, endpoint.Path) {
		return endpoint.Path + param
	}

	return ""
}

func ensureNoTrailingSlash(s string) string {
	return strings.TrimSuffix(s, "/")
}

func ensureLeadingSlash(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}

	return "/" + path
}
