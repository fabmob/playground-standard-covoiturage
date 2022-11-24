# Generated programmatically - DO NOT EDIT

echo "TestGetBookings/test_case"
go run main.go test \
  --method=GET \
  --url=http://localhost:1323/bookings/52fdfc07-2182-654f-163f-5f0f9a621d72 \
  --expectStatus=404

echo "TestGetBookings/test_case#01"
go run main.go test \
  --method=GET \
  --url=http://localhost:1323/bookings/2f8282cb-e2f9-696f-3144-c0aa4ced56db \
  --expectStatus=200

echo "TestGetBookings/test_case#02"
go run main.go test \
  --method=GET \
  --url=http://localhost:1323/bookings/e2807d9c-1dce-26af-00ca-81d4fe11c23e \
  --expectStatus=200

echo "TestPostBookings"
go run main.go test \
  --method=POST \
  --url=http://localhost:1323/bookings \
  --expectStatus=201 \
  <<< '{"driver":{"alias":"","id":"","operator":""},"id":"83472eda-6eb4-7590-6aee-b7f09e757ba9","passenger":{"alias":"","id":"","operator":""},"passengerDropLat":0,"passengerDropLng":0,"passengerPickupDate":0,"passengerPickupLat":0,"passengerPickupLng":0,"price":{},"status":"WAITING_CONFIRMATION"}'

echo "TestPostBookings"
go run main.go test \
  --method=GET \
  --url=http://localhost:1323/bookings/83472eda-6eb4-7590-6aee-b7f09e757ba9 \
  --expectStatus=200

echo "TestPostBookings"
go run main.go test \
  --method=POST \
  --url=http://localhost:1323/bookings \
  --expectStatus=400 \
  <<< '{"driver":{"alias":"","id":"","operator":""},"id":"590c1440-9888-b5b0-7d51-a817ee07c3f2","passenger":{"alias":"","id":"","operator":""},"passengerDropLat":0,"passengerDropLng":0,"passengerPickupDate":0,"passengerPickupLat":0,"passengerPickupLng":0,"price":{},"status":"WAITING_CONFIRMATION"}'

echo "TestPostBookings"
go run main.go test \
  --method=GET \
  --url=http://localhost:1323/bookings/590c1440-9888-b5b0-7d51-a817ee07c3f2 \
  --expectStatus=200

echo "TestPatchBookings"
go run main.go test \
  --method=PATCH \
  --url=http://localhost:1323/bookings/0ad346f9-e692-3ab1-d2f0-91785e9ca0ea?status=VALIDATED \
  --expectStatus=200

echo "TestPatchBookings"
go run main.go test \
  --method=GET \
  --url=http://localhost:1323/bookings/0ad346f9-e692-3ab1-d2f0-91785e9ca0ea \
  --expectStatus=200

echo "TestPatchBookings"
go run main.go test \
  --method=PATCH \
  --url=http://localhost:1323/bookings/68087cc0-282c-35d9-ad8b-51bf6a35a933?status=COMPLETED_PENDING_VALIDATION \
  --expectStatus=200

echo "TestPatchBookings"
go run main.go test \
  --method=GET \
  --url=http://localhost:1323/bookings/68087cc0-282c-35d9-ad8b-51bf6a35a933 \
  --expectStatus=200

echo "TestPatchBookings"
go run main.go test \
  --method=PATCH \
  --url=http://localhost:1323/bookings/3d813194-e9ed-6b09-a1ae-301b83bfdd9d?status=CANCELLED \
  --expectStatus=404

echo "TestPatchBookings"
go run main.go test \
  --method=GET \
  --url=http://localhost:1323/bookings/3d813194-e9ed-6b09-a1ae-301b83bfdd9d \
  --expectStatus=404

echo "TestPatchBookings"
go run main.go test \
  --method=PATCH \
  --url=http://localhost:1323/bookings/1b06f7b5-67c7-f231-9bf3-9f28aa391537?status=VALIDATED \
  --expectStatus=409

echo "TestPatchBookings"
go run main.go test \
  --method=GET \
  --url=http://localhost:1323/bookings/1b06f7b5-67c7-f231-9bf3-9f28aa391537 \
  --expectStatus=200

echo "TestPatchBookings"
go run main.go test \
  --method=PATCH \
  --url=http://localhost:1323/bookings/f84f0c93-2990-ae59-ee94-8e4413ce4e81?status=VALIDATED \
  --expectStatus=409

echo "TestPatchBookings"
go run main.go test \
  --method=GET \
  --url=http://localhost:1323/bookings/f84f0c93-2990-ae59-ee94-8e4413ce4e81 \
  --expectStatus=200

echo "TestPatchBookings"
go run main.go test \
  --method=PATCH \
  --url=http://localhost:1323/bookings/ce140275-2398-b471-e9a9-4ddcec56059b?status=INVALID_STATUS \
  --expectStatus=400

echo "TestPatchBookings"
go run main.go test \
  --method=GET \
  --url=http://localhost:1323/bookings/ce140275-2398-b471-e9a9-4ddcec56059b \
  --expectStatus=200

echo "TestPostBookingEvents"
go run main.go test \
  --method=POST \
  --url=http://localhost:1323/booking_events \
  --expectStatus=200 \
  <<< '{"data":{"id":"6fcf3150-b452-f79a-d30f-524750dbbef4","passengerDropLat":0,"passengerDropLng":0,"passengerPickupDate":0,"passengerPickupLat":0,"passengerPickupLng":0,"status":"WAITING_CONFIRMATION","webUrl":"","driver":{"alias":"","id":"","operator":""},"price":{}},"id":"91523cf5-6600-8472-204b-21603d4a076b","idToken":""}'

echo "TestPostBookingEvents"
go run main.go test \
  --method=GET \
  --url=http://localhost:1323/bookings/6fcf3150-b452-f79a-d30f-524750dbbef4 \
  --expectStatus=200

echo "TestPostBookingEvents"
go run main.go test \
  --method=POST \
  --url=http://localhost:1323/booking_events \
  --expectStatus=200 \
  <<< '{"data":{"id":"cc8c67ad-62d4-b3b1-ee30-02a37a51035f","passengerDropLat":0,"passengerDropLng":0,"passengerPickupDate":0,"passengerPickupLat":0,"passengerPickupLng":0,"status":"CONFIRMED","webUrl":"","driver":{"alias":"","id":"","operator":""},"price":{}},"id":"22128d01-f093-3aca-4106-05310cdc3bb8","idToken":""}'

echo "TestPostBookingEvents"
go run main.go test \
  --method=GET \
  --url=http://localhost:1323/bookings/cc8c67ad-62d4-b3b1-ee30-02a37a51035f \
  --expectStatus=200

echo "TestPostBookingEvents"
go run main.go test \
  --method=POST \
  --url=http://localhost:1323/booking_events \
  --expectStatus=400 \
  <<< '{"data":{"id":"ffda9299-b1d9-fafa-3d47-844c536f73c2","passengerDropLat":0,"passengerDropLng":0,"passengerPickupDate":0,"passengerPickupLat":0,"passengerPickupLng":0,"status":"CONFIRMED","webUrl":"","driver":{"alias":"","id":"","operator":""},"price":{}},"id":"d50fb8fd-a25c-8f1b-114a-976408f9a71b","idToken":""}'

echo "TestPostBookingEvents"
go run main.go test \
  --method=GET \
  --url=http://localhost:1323/bookings/ffda9299-b1d9-fafa-3d47-844c536f73c2 \
  --expectStatus=200

echo "TestPostBookingEvents"
go run main.go test \
  --method=POST \
  --url=http://localhost:1323/booking_events \
  --expectStatus=400 \
  <<< '{"data":{"id":"b2892d57-f402-cd4a-2c11-08cc823ae0c5","passengerDropLat":0,"passengerDropLng":0,"passengerPickupDate":0,"passengerPickupLat":0,"passengerPickupLng":0,"status":"CONFIRMED","webUrl":"","driver":{"alias":"","id":"","operator":""},"price":{}},"id":"90cec22a-723f-cc72-5fb2-462733c2880f","idToken":""}'

echo "TestPostBookingEvents"
go run main.go test \
  --method=GET \
  --url=http://localhost:1323/bookings/b2892d57-f402-cd4a-2c11-08cc823ae0c5 \
  --expectStatus=200

