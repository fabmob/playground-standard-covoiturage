package service

import (
	"fmt"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
)

// keepNFirst keeps n first elements of slice, or returns the slice untouched
// if its length is inferior to n
func keepNFirst[K any](slice []K, n int) []K {
	if len(slice) > n {
		return slice[0:n]
	}

	return slice
}

var statusToIntMap = map[api.BookingStatus]int{
	api.BookingStatusWAITINGCONFIRMATION:        0,
	api.BookingStatusCONFIRMED:                  1,
	api.BookingStatusCOMPLETEDPENDINGVALIDATION: 2,
	api.BookingStatusVALIDATED:                  3,
	api.BookingStatusCANCELLED:                  4,
}

// statusIsAfter checks if status1 is strictly after status2
func statusIsAfter(status1, status2 api.BookingStatus) (bool, error) {
	status1Rank, err := statusRank(status1)
	if err != nil {
		return false, err
	}

	status2Rank, err := statusRank(status2)
	if err != nil {
		return false, err
	}

	return status1Rank > status2Rank, nil
}

func statusRank(status api.BookingStatus) (int, error) {
	statusRank, ok2 := statusToIntMap[status]
	if !ok2 {
		return 0, fmt.Errorf("%s is not a valid status", status)
	}

	return statusRank, nil
}

// errorBody creates an api.BadRequest body from a go error
func errorBody(err error) api.BadRequest {
	errStr := err.Error()
	return api.BadRequest{Error: &errStr}
}
