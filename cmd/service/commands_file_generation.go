package service

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/service/db"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
)

//go:generate bash -c "go test -generate > /dev/null"

var (
	generateTestData bool
	generatedData    = db.NewMockDB()
	commandsFile     = strings.Builder{}

	generatedTestDataFile     = "./db/data/testData.gen.json"
	generatedTestCommandsFile = "../test/commands/testCommands.gen.sh"

	serverEnvVar = "SERVER"
	authEnvVar   = "API_TOKEN"

	unixEpochCounter int64 = 0
	weekInSeconds          = 604800
)

func init() {
	fmt.Fprintln(&commandsFile, "#!/usr/bin/env bash")
	fmt.Fprint(&commandsFile, "# Generated programmatically - DO NOT EDIT\n\n")
	fmt.Fprintf(&commandsFile, "export %s=\"%s\"\n", serverEnvVar, localServer)
	fmt.Fprintf(&commandsFile, "export %s=\"\"\n\n", authEnvVar)
}

// Data needs to be appended once for each test, so we keep track if data has
// already been appended for a given test (with test.Name() as key)
var hasAlreadyAppended = map[string]bool{}

// appendDataIfGenerated is used to populate the `generatedData` db if the
// -generate flag is provided
func appendDataIfGenerated(t *testing.T, mockDB *db.Mock) {
	if _, ok := hasAlreadyAppended[t.Name()]; generateTestData && !ok {
		hasAlreadyAppended[t.Name()] = true
		appendData(mockDB, generatedData)
	}
}

// appendCmdIfGenerated is used to populate the `commands` string, if
// -generate flag is provided
func appendCmdIfGenerated(t *testing.T, request *http.Request, flags test.Flags, body []byte) {
	if generateTestData {
		fmt.Fprint(
			&commandsFile,
			GenerateCommandStr(t, request, flags, body),
		)
	}
}

func appendData(from *db.Mock, to *db.Mock) {
	to.DriverJourneys = append(
		to.GetDriverJourneys(),
		from.GetDriverJourneys()...,
	)

	to.PassengerJourneys = append(
		to.GetPassengerJourneys(),
		from.GetPassengerJourneys()...,
	)

	to.Users = append(
		to.GetUsers(),
		from.GetUsers()...,
	)

	to.DriverRegularTrips = append(
		to.GetDriverRegularTrips(),
		from.GetDriverRegularTrips()...,
	)

	to.PassengerRegularTrips = append(
		to.GetPassengerRegularTrips(),
		from.GetPassengerRegularTrips()...,
	)

	for _, booking := range from.GetBookings() {
		err := to.AddBooking(*booking)
		util.PanicIf(err)
	}
}

// GenerateCommandStr generates a string with the command that should be run
// to test the request with given flags and body.
//
// Used to transform golang tests into command line tests.
func GenerateCommandStr(t *testing.T, request *http.Request, flags test.Flags, body []byte) string {
	var cmd string

	urlWithEnvVar := fmt.Sprintf("$%s%s", serverEnvVar,
		strings.TrimPrefix(request.URL.String(), localServer))

	cmd += fmt.Sprintf("echo \"%s\"\n", t.Name())

	multilineCmd := []string{}
	cmdContinuation := " \\\n  "

	multilineCmd = append(multilineCmd,
		"go run main.go test",
		fmt.Sprintf("--method=%s", request.Method),
		fmt.Sprintf("--url=\"%s\"", urlWithEnvVar),
		fmt.Sprintf("--expectResponseCode=%d", flags.ExpectedResponseCode),
		fmt.Sprintf("--auth=\"$%s\"", authEnvVar),
	)

	if flags.ExpectNonEmpty {
		multilineCmd = append(multilineCmd, "--expectNonEmpty")
	}

	if flags.ExpectedBookingStatus != "" {
		multilineCmd = append(multilineCmd,
			fmt.Sprintf("--expectBookingStatus=%s", flags.ExpectedBookingStatus))
	}

	if body != nil {
		multilineCmd = append(multilineCmd,
			fmt.Sprintf("<<< '%s'", body))
	}

	cmd += strings.Join(multilineCmd, cmdContinuation)
	cmd += "\n\n"
	return cmd
}

// shiftToNextWeek acts as a generator, and yields at each call the
// unix epoch starting at 0.
//
// It is used to isolate as good as possible journeys used for GET
// /driver_journeys, GET /passenger journeys, GET /driver_regular_trip, GET
// /passenger_regular_trip unit tests.
func shiftToNextWeek() {
	unixEpochCounter += int64(weekInSeconds)
}

// setJourneyDatesForGeneration sets, if `generateTestData` == true, a journey date that falls inside the week
// yielded by `shiftToOwnSingleWeek`
func setJourneyDatesForGeneration(schedule *api.JourneySchedule) {
	if generateTestData {
		if schedule.DriverDepartureDate != nil {
			*schedule.DriverDepartureDate += unixEpochCounter
		}
		schedule.PassengerPickupDate += unixEpochCounter
	}
}

func setParamDatesForGeneration(params api.JourneyOrTripPartialParams) {
	switch p := params.(type) {
	case *api.GetDriverJourneysParams:
		p.DepartureDate += int(unixEpochCounter)

	case *api.GetPassengerJourneysParams:
		p.DepartureDate += int(unixEpochCounter)

	case *api.GetDriverRegularTripsParams:
		if p.MinDepartureDate != nil {
			*p.MinDepartureDate += int(unixEpochCounter)
		} else {
			a := int(unixEpochCounter)
			p.MinDepartureDate = &a
		}

		if p.MaxDepartureDate != nil {
			*p.MaxDepartureDate += int(unixEpochCounter) + weekInSeconds
		} else {
			a := int(unixEpochCounter) + weekInSeconds
			p.MaxDepartureDate = &a
		}

	case *api.GetPassengerRegularTripsParams:
		if p.MinDepartureDate != nil {
			*p.MinDepartureDate += int(unixEpochCounter)
		} else {
			a := int(unixEpochCounter)
			p.MinDepartureDate = &a
		}

		if p.MaxDepartureDate != nil {
			*p.MaxDepartureDate += int(unixEpochCounter) + weekInSeconds
		} else {
			a := int(unixEpochCounter) + weekInSeconds
			p.MaxDepartureDate = &a
		}
	}
}
