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
[binary](https://github.com/fabmob/playground-standard-covoiturage/blob/documentation/pscovoit) 
for linux (may not be the latest development version), and make it executable, 
or clone the repo and enter `go build -o pscovoit` in the root folder.



## Run the fake server

The `serve` subcommand runs the server on https://localhost:1323 (port not 
customizable yet):

```
./pscovoit serve
```

