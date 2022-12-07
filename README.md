# Playground standard covoiturage

A tool to test an API against the [standard-covoiturage 
specification](https://github.com/fabmob/standard-covoiturage):
- A test server with fake in-memory data to play with the standard API.
- A client that runs a test against a server request.

The tool is currently being developped. More about the aimed functional scope 
[here (fr)](./docs/proposition_fonctionelle.pdf).

## Install

You can find and download the executable for various platforms 
[here](https://github.com/fabmob/playground-standard-covoiturage/releases), or 
clone the repo and enter `go build -o pscovoit` in the root folder.

## Run the fake server

The `serve` subcommand runs the server on http://localhost:1323 (port not 
customizable yet):

```sh
./pscovoit serve
```

The served data can be inspected 
[here](https://github.com/fabmob/playground-standard-covoiturage/blob/main/cmd/service/data/defaultData.json), 
or custom data can be used with the `--data` flag pointing to a valid json 
data file (check-out [type 
`MockDBDataInterface`](https://github.com/fabmob/playground-standard-covoiturage/blob/eb4ccb0cb125639921394f851a7e975e07cbc386/cmd/service/db/db.go#L127) 
for more details on data structure). 

## Test a request

The `test` subcommand runs tests on a given request. 

There are several ways to specify the request, the three examples below are 
equivalent:

- through an url (̀`method` flag is set to GET by default and could be omitted 
  here)

```sh
./pscovoit test \
--url "http://localhost:1323/driver_journeys?arrivalLat=48.8450234&arrivalLng=2.3997529&departureDate=1665579951&departureLat=47.461737&departureLng=1.061393" \
--method=GET
```

- through an url with (some or all) query parameters as flags

```sh
./pscovoit test --url "http://localhost:1323/driver_journeys" \
  -q arrivalLat=48.8450234 \
  -q arrivalLng=2.3997529 \
  -q departureDate=1665579951 \
  -q departureLat=47.461737 \
  -q departureLng=1.061393
```

- through subcommands 
  
```sh
./pscovoit test get driverJourneys \
  --server "http://localhost:1323" \
  --arrivalLat=48.8450234 \
  --arrivalLng=2.3997529 \
  --departureDate=1665579951 \
  --departureLat=47.461737 \
  --departureLng=1.061393
```

A request body can be read from standard input, e.g.:
 
- Read from file

```sh
pscovoit test post bookings <file_with_body.txt
```

- Passed as string

```sh
pscovoit test post bookings <<< "{body}"
```

By default, only the failed tests are reported. Use the `--verbose` flag to 
see all tests and additional information.


## Autocompletion

The last method may greatly benefit from autocompletion.

This is possible if the binary is in $PATH, see for more information:
```
./pscovoit completion --help
```

## Use in CI

The tool returns exit code 1 in case of an assertion failure. See a simple 
example of use in github CI 
[here](https://github.com/fabmob/playground-standard-covoiturage/blob/eb4ccb0cb125639921394f851a7e975e07cbc386/.github/workflows/go_tests_lint.yml#L42) 
(github workflow) and 
[here](https://github.com/fabmob/playground-standard-covoiturage/tree/main/cmd/test/commands) 
(test commands scripts). 

## Tests and assertions

### Test Flags

* `--expectResponseCode`: additional check that the HTTP response code is as 
  expected
* `--expectNonEmpty` (array responses): additionnal check that this array is 
  not empty.
* `--expectBookingStatus` (GET /bookings): additional check that the booking 
  has the expected booking status. 
  
### Example tests

For an example of a thorough test suite for an API, look at
[this example](./cmd/service/data/testCommands.gen.sh),
that works together with [this data](./cmd/service/data/testData.gen.json).


### GET /driver_journey, GET /passenger_journey, GET /driver_regular_trips, GET/passenger_regular_trips

The following assertions are run on these two endpoints :
- assert format
- assert response status code 200 (optional)
- assert header Content-Type:application/json
- assert query parameter "departureRadius" 
- assert query parameter "arrivalRadius"
- assert query parameter "timeDelta"
- assert query parameter "count"
- assert unique ids
- assert response property "operator"

### POST /bookings, POST /booking_events, PATCH /bookings, POST /messages

- assert format
- assert response status code (optional)

### GET /bookings
 
- assert format
- assert response status code (optional)
- assert booking status (optional)

### Assertions reference

In alphabetic order: 

| Assertion code                 | description                                                                                                                                            |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| assert API call success        | Checks that the response data has been succesfully collected                                                                                           |
| assert format                  | Checks that the format of the response complies to the standard's openAPI specification. Especially, the observed status code needs to be documented.  |
| assert header X:Y              | Checks that the response has header X with value Y.                                                                                                    |
| assert query parameter X       | Checks that the response complies to the expectations of the queryparameter X.                                                                         |
| assert response not empty      | Checks that the response is not an empty array.                                                                                                        |
| assert response property X     | Checks that the response property X meets the expectations given by the standard.                                                                      |
| assert response status code X  | Checks that the status code X is returned.                                                                                                             |
| assert unique ids              | Checks that the response objects have no duplicated "id" property.                                                                                     |


## Release

Releases are made with [goreleaser](https://goreleaser.com/quick-start/). 

## Roadmap until first release

* [ ] Fix issue #25
* [ ] Implement server with unit tests
    * [X] All but regular trips
* [ ] Export unit test data and generate make unit test commands for reuse
    * [X] All but search Endpoints
* [X] Cross platform release with goreleaser
* [ ] Add new custom assertions to current tests
* [X] Update documentation
* [X] Add license
* [X] Add Changelog (auto generated with release)
* [ ] Tag and release in github 
