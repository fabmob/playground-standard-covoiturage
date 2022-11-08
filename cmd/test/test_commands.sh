#!/usr/bin/env bash

# Testing different testing commands. To be run inside the project root
# directory. The commands are not expected to test any data.

# exit when any command fails
set -e
set -o pipefail

# Clean up subprocesses on exit
# Do not use builtin bash `kill` command
enable -n kill
trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

echo "Run fake server"
go run main.go serve > /dev/null &
sleep 1

echo "Test GET /driver_journeys with url"
go run main.go test \
  --url "http://localhost:1323/driver_journeys" \
  -q departureLat=0 \
  -q departureLng=0 \
  -q arrivalLat=0 \
  -q arrivalLng=0 \
  -q departureDate=0 \
  -q timeDelta=0 \
  -q departureRadius=0 \
  -q arrivalRadius=0


echo "Test GET /passenger_journeys with short command, no optional query parameter"
go run main.go test get passengerJourneys --server "http://localhost:1323" --departureLat=0 --departureLng=0 --arrivalLat=0 --arrivalLng=0 --departureDate=0


echo "Test GET /passenger_journeys with short command, all optional query parameter"
go run main.go test get passengerJourneys --server "http://localhost:1323" --departureLat=0 --departureLng=0 --arrivalLat=0 --arrivalLng=0 --departureDate=0

# echo "Test GET /bookings/{bookingId} with short command"
# go run main.go test get bookings --server "http://localhost:1323" --bookingId="42"

