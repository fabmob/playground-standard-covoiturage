# Playground standard covoiturage

A tool to test an API against the [standard-covoiturage 
specification](https://github.com/fabmob/standard-covoiturage):
- A test server with fake in-memory data to play with the standard API.
- A client that runs a test against a server request.

The tool is currently being developped. More about the aimed functional scope 
[here (fr)](./docs/proposition_fonctionelle.pdf).

## Install

No fancy installation mechanism available yet.

Download the 
[binary](https://github.com/fabmob/playground-standard-covoiturage/blob/main/pscovoit) 
for linux (may not be the latest development version), and make it executable, 
or clone the repo and enter `go build -o pscovoit` in the root folder.



## Run the fake server

The `serve` subcommand runs the server on http://localhost:1323 (port not 
customizable yet):

```sh
./pscovoit serve
```

The served data can be inspected 
[here](https://github.com/fabmob/playground-standard-covoiturage/blob/main/cmd/service/data/defaultData.json) 
and is not yet customizable. 

## Test a request

The `test` subcommand runs tests on a given request. 

There are several ways to specify the request, the three examples below are 
equivalent:

- through an url

```sh
./pscovoit test --url 
"http://localhost:1323/driver_journeys?arrivalLat=48.8450234&arrivalLng=2.3997529&departureDate=1665579951&departureLat=47.461737&departureLng=1.061393"
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
  --server "http://localhost:1323"
  --arrivalLat=48.8450234 \
  --arrivalLng=2.3997529 \
  --departureDate=1665579951 \
  --departureLat=47.461737 \
  --departureLng=1.061393
```

By default, only the failed tests are reported. Use the `--verbose` flag to 
see all tests. 

## Autocompletion

The last method may greatly benefit from autocompletion.

This is possible if the binary is in $PATH, see for more information:
```
./pscovoit completion --help
```

## Tests and assertions

### Test for non-empty response

The `--disallowEmpty` flag runs for all endpoints that returns an array an 
additionnal check that this array is not empty. 

### GET /driver_journey and GET /passenger_journey

The following assertions are run on these two endpoints :
- assert response status code 200
- assert header Content-Type:application/json
- assert format
- assert query parameter "departureRadius"
- assert query parameter "arrivalRadius"
- assert query parameter "timeDelta"
- assert query parameter "count"
- assert unique ids
- assert response property "operator"

See below for assertions reference. 

### Assertions reference

### assert API call success

Checks that the response data has been succesfully collected

### assert response not empty

Checks that the response is not an empty array.

### assert response status code X

Checks that the status code X is returned.

### assert header X:Y

Checks that the response has header X with value Y.

### assert format

Checks that the format of the response complies to the standards openAPI 
specification.

### assert query parameter X

Checks that the response complies to the expectations of the queryparameter X.

### assert unique ids

Checks that the response objects have no duplicated "id" field.

### assert response property X

Checks that the response property X meets the expectations given by the 
standard.
