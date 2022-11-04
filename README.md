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

The `serve` subcommand runs the server on https://localhost:1323 (port not 
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
"https://localhost:1373/driver_journeys?arrivalLat=48.8450234&arrivalLng=2.3997529&departureDate=1665579951&departureLat=47.461737&departureLng=1.061393"
```

- through an url with (some or all) query parameters as flags

```sh
./pscovoit test --url "https://localhost:1373/driver_journeys" \
  -q arrivalLat=48.8450234 \
  -q arrivalLng=2.3997529 \
  -q departureDate=1665579951 \
  -q departureLat=47.461737 \
  -q departureLng=1.061393
```

- through subcommands 
  
```sh
./pscovoit test get driverJourneys \
  --server "https://localhost:1373"
  --arrivalLat=48.8450234 \
  --arrivalLng=2.3997529 \
  --departureDate=1665579951 \
  --departureLat=47.461737 \
  --departureLng=1.061393
```

## Autocompletion

