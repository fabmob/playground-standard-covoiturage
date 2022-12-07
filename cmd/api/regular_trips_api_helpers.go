package api

type DriverRegularTrip struct {
	DriverTrip
	Schedules *[]Schedule
}

type PassengerRegularTrip struct {
	PassengerTrip
	Schedules *[]Schedule
}
