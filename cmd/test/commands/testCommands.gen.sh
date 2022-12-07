#!/usr/bin/env bash
# Generated programmatically - DO NOT EDIT

export SERVER="http://localhost:1323"
export API_TOKEN=""

echo "TestDriverJourneys/No_data"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=604800&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestDriverJourneys/Departure_radius_1"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=1209600&departureLat=46.160454&departureLng=-1.2219607&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestDriverJourneys/Departure_radius_2"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=1814400&departureLat=46.160454&departureLng=-1.2219607&departureRadius=2&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestDriverJourneys/Departure_radius_3"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=2419200&departureLat=46.160454&departureLng=-1.2219607&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestDriverJourneys/Departure_radius_3#01"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=3024000&departureLat=46.160454&departureLng=-1.2219607&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestDriverJourneys/Arrival_radius_1"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=1&departureDate=3628800&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestDriverJourneys/Arrival_radius_2"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=2&departureDate=4233600&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestDriverJourneys/Arrival_radius_3"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=1&departureDate=4838400&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestDriverJourneys/Arrival_radius_4"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=1&departureDate=5443200&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestDriverJourneys/Count_1"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=1&departureDate=6048000&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestDriverJourneys/Count_2"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=0&departureDate=6652800&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestDriverJourneys/Count_3"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=2&departureDate=7257600&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestDriverJourneys/Count_4_-_count_>_n_driver_journeys"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=1&departureDate=7862400&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestDriverJourneys/TimeDelta_1"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=8467200&departureLat=0&departureLng=0&departureRadius=1&timeDelta=10" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestDriverJourneys/TimeDelta_2"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=9072000&departureLat=0&departureLng=0&departureRadius=1&timeDelta=10" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestDriverJourneys/TimeDelta_3"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=9676800&departureLat=0&departureLng=0&departureRadius=1&timeDelta=20" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestPassengerJourneys/No_data"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=10281600&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestPassengerJourneys/Departure_radius_1"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=10886400&departureLat=46.160454&departureLng=-1.2219607&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestPassengerJourneys/Departure_radius_2"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=11491200&departureLat=46.160454&departureLng=-1.2219607&departureRadius=2&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestPassengerJourneys/Departure_radius_3"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=12096000&departureLat=46.160454&departureLng=-1.2219607&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestPassengerJourneys/Departure_radius_3#01"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=12700800&departureLat=46.160454&departureLng=-1.2219607&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestPassengerJourneys/Arrival_radius_1"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=1&departureDate=13305600&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestPassengerJourneys/Arrival_radius_2"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=2&departureDate=13910400&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestPassengerJourneys/Arrival_radius_3"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=1&departureDate=14515200&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestPassengerJourneys/Arrival_radius_4"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=1&departureDate=15120000&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestPassengerJourneys/Count_1"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=1&departureDate=15724800&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestPassengerJourneys/Count_2"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=0&departureDate=16329600&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestPassengerJourneys/Count_3"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=2&departureDate=16934400&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestPassengerJourneys/Count_4_-_count_>_n_driver_journeys"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=1&departureDate=17539200&departureLat=0&departureLng=0&departureRadius=1&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestPassengerJourneys/TimeDelta_1"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=18144000&departureLat=0&departureLng=0&departureRadius=1&timeDelta=10" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestPassengerJourneys/TimeDelta_2"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=18748800&departureLat=0&departureLng=0&departureRadius=1&timeDelta=10" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestPassengerJourneys/TimeDelta_3"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_journeys?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureDate=19353600&departureLat=0&departureLng=0&departureRadius=1&timeDelta=20" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetDriverRegularTrips/No_data"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=20563200&minDepartureDate=19958400&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestGetDriverRegularTrips/Departure_radius_1"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureLat=46.160454&departureLng=-1.2219607&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=21168000&minDepartureDate=20563200&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetDriverRegularTrips/Departure_radius_2"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureLat=46.160454&departureLng=-1.2219607&departureRadius=2&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=21772800&minDepartureDate=21168000&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetDriverRegularTrips/Departure_radius_3"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureLat=46.160454&departureLng=-1.2219607&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=22377600&minDepartureDate=21772800&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestGetDriverRegularTrips/Departure_radius_3#01"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureLat=46.160454&departureLng=-1.2219607&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=22982400&minDepartureDate=22377600&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetDriverRegularTrips/Arrival_radius_1"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_regular_trips?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=1&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=23587200&minDepartureDate=22982400&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetDriverRegularTrips/Arrival_radius_2"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_regular_trips?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=2&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=24192000&minDepartureDate=23587200&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetDriverRegularTrips/Arrival_radius_3"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_regular_trips?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=1&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=24796800&minDepartureDate=24192000&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestGetDriverRegularTrips/Arrival_radius_4"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_regular_trips?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=1&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=25401600&minDepartureDate=24796800&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetDriverRegularTrips/Count_1"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=1&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=26006400&minDepartureDate=25401600&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetDriverRegularTrips/Count_2"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=0&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=26611200&minDepartureDate=26006400&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestGetDriverRegularTrips/Count_3"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=2&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=27216000&minDepartureDate=26611200&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetDriverRegularTrips/Count_4_-_count_>_n_driver_journeys"
go run main.go test \
  --method=GET \
  --url="$SERVER/driver_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=1&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=27820800&minDepartureDate=27216000&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestGetPassengerRegularTrips/No_data"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=28425600&minDepartureDate=27820800&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestGetPassengerRegularTrips/Departure_radius_1"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureLat=46.160454&departureLng=-1.2219607&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=29030400&minDepartureDate=28425600&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetPassengerRegularTrips/Departure_radius_2"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureLat=46.160454&departureLng=-1.2219607&departureRadius=2&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=29635200&minDepartureDate=29030400&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetPassengerRegularTrips/Departure_radius_3"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureLat=46.160454&departureLng=-1.2219607&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=30240000&minDepartureDate=29635200&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestGetPassengerRegularTrips/Departure_radius_3#01"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&departureLat=46.160454&departureLng=-1.2219607&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=30844800&minDepartureDate=30240000&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetPassengerRegularTrips/Arrival_radius_1"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_regular_trips?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=1&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=31449600&minDepartureDate=30844800&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetPassengerRegularTrips/Arrival_radius_2"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_regular_trips?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=2&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=32054400&minDepartureDate=31449600&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetPassengerRegularTrips/Arrival_radius_3"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_regular_trips?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=1&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=32659200&minDepartureDate=32054400&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestGetPassengerRegularTrips/Arrival_radius_4"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_regular_trips?arrivalLat=46.160454&arrivalLng=-1.2219607&arrivalRadius=1&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=33264000&minDepartureDate=32659200&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetPassengerRegularTrips/Count_1"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=1&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=33868800&minDepartureDate=33264000&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetPassengerRegularTrips/Count_2"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=0&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=34473600&minDepartureDate=33868800&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestGetPassengerRegularTrips/Count_3"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=2&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=35078400&minDepartureDate=34473600&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetPassengerRegularTrips/Count_4_-_count_>_n_driver_journeys"
go run main.go test \
  --method=GET \
  --url="$SERVER/passenger_regular_trips?arrivalLat=0&arrivalLng=0&arrivalRadius=1&count=1&departureLat=0&departureLng=0&departureRadius=1&departureTimeOfDay=08%3A00%3A00&maxDepartureDate=35683200&minDepartureDate=35078400&timeDelta=900" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestGetBookings/getting_a_non-existing_booking_returns_code_404"
go run main.go test \
  --method=GET \
  --url="$SERVER/bookings/52fdfc07-2182-654f-163f-5f0f9a621d72" \
  --expectResponseCode=404 \
  --auth="$API_TOKEN"

echo "TestGetBookings/getting_an_existing_booking_returns_it_with_code_200_#1"
go run main.go test \
  --method=GET \
  --url="$SERVER/bookings/2f8282cb-e2f9-696f-3144-c0aa4ced56db" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestGetBookings/getting_an_existing_booking_returns_it_with_code_200_#2"
go run main.go test \
  --method=GET \
  --url="$SERVER/bookings/e2807d9c-1dce-26af-00ca-81d4fe11c23e" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestPostBookings/Posting_a_new_booking_succeeds_with_code_201"
go run main.go test \
  --method=POST \
  --url="$SERVER/bookings" \
  --expectResponseCode=201 \
  --auth="$API_TOKEN" \
  <<< '{"driver":{"alias":"","id":"","operator":""},"id":"83472eda-6eb4-7590-6aee-b7f09e757ba9","passenger":{"alias":"","id":"","operator":""},"passengerDropLat":0,"passengerDropLng":0,"passengerPickupDate":0,"passengerPickupLat":0,"passengerPickupLng":0,"price":{},"status":"WAITING_CONFIRMATION"}'

echo "TestPostBookings/Posting_a_new_booking_succeeds_with_code_201"
go run main.go test \
  --method=GET \
  --url="$SERVER/bookings/83472eda-6eb4-7590-6aee-b7f09e757ba9" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestPostBookings/Posting_a_booking_with_colliding_ID_fails_with_code_400"
go run main.go test \
  --method=POST \
  --url="$SERVER/bookings" \
  --expectResponseCode=400 \
  --auth="$API_TOKEN" \
  <<< '{"driver":{"alias":"","id":"","operator":""},"id":"590c1440-9888-b5b0-7d51-a817ee07c3f2","passenger":{"alias":"","id":"","operator":""},"passengerDropLat":0,"passengerDropLng":0,"passengerPickupDate":0,"passengerPickupLat":0,"passengerPickupLng":0,"price":{},"status":"WAITING_CONFIRMATION"}'

echo "TestPostBookings/Posting_a_booking_with_colliding_ID_fails_with_code_400"
go run main.go test \
  --method=GET \
  --url="$SERVER/bookings/590c1440-9888-b5b0-7d51-a817ee07c3f2" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectNonEmpty

echo "TestPatchBookings/patching_VALIDATED_over_WAITING_CONFIRMATION_succeeds"
go run main.go test \
  --method=PATCH \
  --url="$SERVER/bookings/0ad346f9-e692-3ab1-d2f0-91785e9ca0ea?status=VALIDATED" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestPatchBookings/patching_VALIDATED_over_WAITING_CONFIRMATION_succeeds"
go run main.go test \
  --method=GET \
  --url="$SERVER/bookings/0ad346f9-e692-3ab1-d2f0-91785e9ca0ea" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectBookingStatus=VALIDATED

echo "TestPatchBookings/patching_COMPLETED_PENDING_VALIDATION_over_WAITING_CONFIRMATION_succeeds"
go run main.go test \
  --method=PATCH \
  --url="$SERVER/bookings/68087cc0-282c-35d9-ad8b-51bf6a35a933?status=COMPLETED_PENDING_VALIDATION" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN"

echo "TestPatchBookings/patching_COMPLETED_PENDING_VALIDATION_over_WAITING_CONFIRMATION_succeeds"
go run main.go test \
  --method=GET \
  --url="$SERVER/bookings/68087cc0-282c-35d9-ad8b-51bf6a35a933" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectBookingStatus=COMPLETED_PENDING_VALIDATION

echo "TestPatchBookings/patching_a_non-existing_booking_returns_code_404"
go run main.go test \
  --method=PATCH \
  --url="$SERVER/bookings/3d813194-e9ed-6b09-a1ae-301b83bfdd9d?status=CANCELLED" \
  --expectResponseCode=404 \
  --auth="$API_TOKEN"

echo "TestPatchBookings/patching_a_non-existing_booking_returns_code_404"
go run main.go test \
  --method=GET \
  --url="$SERVER/bookings/3d813194-e9ed-6b09-a1ae-301b83bfdd9d" \
  --expectResponseCode=404 \
  --auth="$API_TOKEN"

echo "TestPatchBookings/patching_VALIDATED_other_VALIDATED_fails_with_code_409"
go run main.go test \
  --method=PATCH \
  --url="$SERVER/bookings/1b06f7b5-67c7-f231-9bf3-9f28aa391537?status=VALIDATED" \
  --expectResponseCode=409 \
  --auth="$API_TOKEN"

echo "TestPatchBookings/patching_VALIDATED_other_VALIDATED_fails_with_code_409"
go run main.go test \
  --method=GET \
  --url="$SERVER/bookings/1b06f7b5-67c7-f231-9bf3-9f28aa391537" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectBookingStatus=VALIDATED

echo "TestPatchBookings/patching_VALIDATED_other_CANCELLED_fails_with_code_409"
go run main.go test \
  --method=PATCH \
  --url="$SERVER/bookings/f84f0c93-2990-ae59-ee94-8e4413ce4e81?status=VALIDATED" \
  --expectResponseCode=409 \
  --auth="$API_TOKEN"

echo "TestPatchBookings/patching_VALIDATED_other_CANCELLED_fails_with_code_409"
go run main.go test \
  --method=GET \
  --url="$SERVER/bookings/f84f0c93-2990-ae59-ee94-8e4413ce4e81" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectBookingStatus=CANCELLED

echo "TestPatchBookings/patching_INVALID_STATUS_fails_with_code_400"
go run main.go test \
  --method=PATCH \
  --url="$SERVER/bookings/ce140275-2398-b471-e9a9-4ddcec56059b?status=INVALID_STATUS" \
  --expectResponseCode=400 \
  --auth="$API_TOKEN"

echo "TestPatchBookings/patching_INVALID_STATUS_fails_with_code_400"
go run main.go test \
  --method=GET \
  --url="$SERVER/bookings/ce140275-2398-b471-e9a9-4ddcec56059b" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectBookingStatus=WAITING_CONFIRMATION

echo "TestPostBookingEvents/posting_a_new_bookingEvent_with_status_WAITING_CONFIRMATION_succeeds"
go run main.go test \
  --method=POST \
  --url="$SERVER/booking_events" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  <<< '{"data":{"id":"6fcf3150-b452-f79a-d30f-524750dbbef4","passengerDropLat":0,"passengerDropLng":0,"passengerPickupDate":0,"passengerPickupLat":0,"passengerPickupLng":0,"status":"WAITING_CONFIRMATION","webUrl":"","driver":{"alias":"","id":"","operator":""},"price":{}},"id":"91523cf5-6600-8472-204b-21603d4a076b","idToken":""}'

echo "TestPostBookingEvents/posting_a_new_bookingEvent_with_status_WAITING_CONFIRMATION_succeeds"
go run main.go test \
  --method=GET \
  --url="$SERVER/bookings/6fcf3150-b452-f79a-d30f-524750dbbef4" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectBookingStatus=WAITING_CONFIRMATION

echo "TestPostBookingEvents/posting_a_bookingEvent_on_existing_booking_(status_CONFIRMED_over_WAITING_CONFIRMATION)_changes_its_status"
go run main.go test \
  --method=POST \
  --url="$SERVER/booking_events" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  <<< '{"data":{"id":"cc8c67ad-62d4-b3b1-ee30-02a37a51035f","passengerDropLat":0,"passengerDropLng":0,"passengerPickupDate":0,"passengerPickupLat":0,"passengerPickupLng":0,"status":"CONFIRMED","webUrl":"","driver":{"alias":"","id":"","operator":""},"price":{}},"id":"22128d01-f093-3aca-4106-05310cdc3bb8","idToken":""}'

echo "TestPostBookingEvents/posting_a_bookingEvent_on_existing_booking_(status_CONFIRMED_over_WAITING_CONFIRMATION)_changes_its_status"
go run main.go test \
  --method=GET \
  --url="$SERVER/bookings/cc8c67ad-62d4-b3b1-ee30-02a37a51035f" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectBookingStatus=CONFIRMED

echo "TestPostBookingEvents/posting_a_bookingEvent_on_existing_booking_(status_CONFIRMED_over_CONFIRMED)_fails_with_code_400"
go run main.go test \
  --method=POST \
  --url="$SERVER/booking_events" \
  --expectResponseCode=400 \
  --auth="$API_TOKEN" \
  <<< '{"data":{"id":"ffda9299-b1d9-fafa-3d47-844c536f73c2","passengerDropLat":0,"passengerDropLng":0,"passengerPickupDate":0,"passengerPickupLat":0,"passengerPickupLng":0,"status":"CONFIRMED","webUrl":"","driver":{"alias":"","id":"","operator":""},"price":{}},"id":"d50fb8fd-a25c-8f1b-114a-976408f9a71b","idToken":""}'

echo "TestPostBookingEvents/posting_a_bookingEvent_on_existing_booking_(status_CONFIRMED_over_CONFIRMED)_fails_with_code_400"
go run main.go test \
  --method=GET \
  --url="$SERVER/bookings/ffda9299-b1d9-fafa-3d47-844c536f73c2" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectBookingStatus=CONFIRMED

echo "TestPostBookingEvents/posting_a_bookingEvent_on_existing_booking_(status_CONFIRMED_over_CANCELLED)_fails_with_code_400"
go run main.go test \
  --method=POST \
  --url="$SERVER/booking_events" \
  --expectResponseCode=400 \
  --auth="$API_TOKEN" \
  <<< '{"data":{"id":"b2892d57-f402-cd4a-2c11-08cc823ae0c5","passengerDropLat":0,"passengerDropLng":0,"passengerPickupDate":0,"passengerPickupLat":0,"passengerPickupLng":0,"status":"CONFIRMED","webUrl":"","driver":{"alias":"","id":"","operator":""},"price":{}},"id":"90cec22a-723f-cc72-5fb2-462733c2880f","idToken":""}'

echo "TestPostBookingEvents/posting_a_bookingEvent_on_existing_booking_(status_CONFIRMED_over_CANCELLED)_fails_with_code_400"
go run main.go test \
  --method=GET \
  --url="$SERVER/bookings/b2892d57-f402-cd4a-2c11-08cc823ae0c5" \
  --expectResponseCode=200 \
  --auth="$API_TOKEN" \
  --expectBookingStatus=CANCELLED

echo "TestPostMessage/Posting_message_with_both_user_known_succeeds_with_code_201"
go run main.go test \
  --method=POST \
  --url="$SERVER/messages" \
  --expectResponseCode=201 \
  --auth="$API_TOKEN" \
  <<< '{"from":{"alias":"alice","id":"2","operator":"default.operator.com"},"message":"some message","recipientCarpoolerType":"DRIVER","to":{"alias":"bob","id":"1","operator":"default.operator.com"}}'

echo "TestPostMessage/Posting_message_with_recipient_unknown_fails_with_code_404"
go run main.go test \
  --method=POST \
  --url="$SERVER/messages" \
  --expectResponseCode=404 \
  --auth="$API_TOKEN" \
  <<< '{"from":{"alias":"carole","id":"3","operator":"default.operator.com"},"message":"some message","recipientCarpoolerType":"DRIVER","to":{"alias":"david","id":"4","operator":"default.operator.com"}}'

echo "TestPostMessage/Posting_message_with_sender_unknown_succeeds_with_code_201"
go run main.go test \
  --method=POST \
  --url="$SERVER/messages" \
  --expectResponseCode=201 \
  --auth="$API_TOKEN" \
  <<< '{"from":{"alias":"eve","id":"5","operator":"default.operator.com"},"message":"some message","recipientCarpoolerType":"DRIVER","to":{"alias":"fanny","id":"6","operator":"default.operator.com"}}'

