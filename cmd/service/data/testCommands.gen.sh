# Generated programmatically - DO NOT EDIT

# TestGetBookings/test_case
go run main.go test \
  --url=http://localhost:1323/bookings/52fdfc07-2182-654f-163f-5f0f9a621d72 \
  --expectStatus=404

# TestGetBookings/test_case#01
go run main.go test \
  --url=http://localhost:1323/bookings/2f8282cb-e2f9-696f-3144-c0aa4ced56db \
  --expectStatus=200

# TestGetBookings/test_case#02
go run main.go test \
  --url=http://localhost:1323/bookings/e2807d9c-1dce-26af-00ca-81d4fe11c23e \
  --expectStatus=200

# # TestPostBookings
# go run main.go test \
#   --url=http://localhost:1323/bookings/83472eda-6eb4-7590-6aee-b7f09e757ba9 \
#   --expectStatus=200

# # TestPostBookings
# go run main.go test \
#   --url=http://localhost:1323/bookings/590c1440-9888-b5b0-7d51-a817ee07c3f2 \
#   --expectStatus=200

# # TestPatchBookings
# go run main.go test \
#   --url=http://localhost:1323/bookings/0ad346f9-e692-3ab1-d2f0-91785e9ca0ea \
#   --expectStatus=200

# # TestPatchBookings
# go run main.go test \
#   --url=http://localhost:1323/bookings/68087cc0-282c-35d9-ad8b-51bf6a35a933 \
#   --expectStatus=200

# # TestPatchBookings
# go run main.go test \
#   --url=http://localhost:1323/bookings/3d813194-e9ed-6b09-a1ae-301b83bfdd9d \
#   --expectStatus=404

# # TestPatchBookings
# go run main.go test \
#   --url=http://localhost:1323/bookings/1b06f7b5-67c7-f231-9bf3-9f28aa391537 \
#   --expectStatus=200

# # TestPatchBookings
# go run main.go test \
#   --url=http://localhost:1323/bookings/f84f0c93-2990-ae59-ee94-8e4413ce4e81 \
#   --expectStatus=200

# # TestPatchBookings
# go run main.go test \
#   --url=http://localhost:1323/bookings/ce140275-2398-b471-e9a9-4ddcec56059b \
#   --expectStatus=200

# # TestPostBookingEvents
# go run main.go test \
#   --url=http://localhost:1323/bookings/6fcf3150-b452-f79a-d30f-524750dbbef4 \
#   --expectStatus=200

# # TestPostBookingEvents
# go run main.go test \
#   --url=http://localhost:1323/bookings/cc8c67ad-62d4-b3b1-ee30-02a37a51035f \
#   --expectStatus=200

# # TestPostBookingEvents
# go run main.go test \
#   --url=http://localhost:1323/bookings/ffda9299-b1d9-fafa-3d47-844c536f73c2 \
#   --expectStatus=200

# # TestPostBookingEvents
# go run main.go test \
#   --url=http://localhost:1323/bookings/b2892d57-f402-cd4a-2c11-08cc823ae0c5 \
#   --expectStatus=200

