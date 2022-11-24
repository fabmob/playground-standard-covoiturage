#!/usr/bin/env bash

# Testing different testing commands. To be run inside the project root
# directory. The commands are not expected to test any data.

# The fake server must run on http://localhost:1323 for the tests to pass

# exit when any command fails
set -e
set -o pipefail


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

echo "Test GET /bookings/{bookingId} with short command"
go run main.go test get bookings \
  --server "http://localhost:1323" \
  --bookingId="123e4567-e89b-12d3-a456-426614174000" \
  --expectStatus=404

echo "Test POST /bookings with body"
go run main.go test post bookings \
  --server "http://localhost:1323" \
  <<< '{
  "driver": {
    "alias": "abc87",
    "id": "12345"
  },
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "passenger": {
    "alias": "cde69",
    "id": "67890"
  },
  "passengerDropLat": 45.8275,
  "passengerDropLng":  1.25987,
  "status": "WAITING_CONFIRMATION"
}'
