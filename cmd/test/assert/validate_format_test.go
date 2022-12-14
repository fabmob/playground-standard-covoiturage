package assert

import (
	"context"
	"net/http"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
)

func TestUndocumentedStatusCode(t *testing.T) {
	var (
		invalidStatus = http.StatusLoopDetected
	)
	request, err := api.NewGetDriverJourneysRequest(localServer, &api.GetDriverJourneysParams{})
	util.PanicIf(err)

	response := mockResponse(invalidStatus, "[]", nil)

	validationErr := validateResponse(request, response)
	if validationErr == nil {
		t.Error("Format validation is expected to fail for undocumented status code")
	}

}

func TestFindRoute(t *testing.T) {
	server := endpoint.Server("https://abc.fr/abc")
	url := string(server) + "/driver_journeys"

	request, err := http.NewRequest("GET", url, nil)
	util.PanicIf(err)

	ctx := context.Background()

	_, _, err = findRoute(ctx, request, server)
	if err != nil {
		t.Log(server)
		t.Error("Format validation with kin-openapi does not find route properly")
	}
}
