# utm-fwd

utm-fwd is a Go WGS84 UTM forward-projection calculator. It computes the UTM
zone and central meridian from longitude, then projects latitude/longitude to
easting and northing with a transverse Mercator series on the WGS84 ellipsoid.
The result includes point scale `k`, grid convergence, latitude band, and an
explicit flag when a caller forces a zone different from the automatic zone.
The same calculations are exposed through HTTP JSON endpoints and CLI
subcommands with no web page.

## Usage

Run the HTTP server:

```bash
go run . serve -addr :8080
```

Evaluate from the command line:

```bash
go run . forward -lat 39.9 -lon 116.4
go run . forward -lat 39.9 -lon 116.4 -zone 51
go run . zone -lon 116.4
```

Run the Beijing example:

```bash
go run . example -file example/beijing.json
```

Beijing falls in zone 50 with central meridian 117°E. Its easting is near
449 km, inside the expected 400-600 km range.

## HTTP API

```text
POST /api/forward  {"lat":39.9,"lon":116.4}
POST /api/forward  {"lat":39.9,"lon":116.4,"zone":51}
POST /api/zone     {"lon":116.4}
GET  /health
```

Latitudes outside `[-80,84]`, longitudes outside `[-180,180]`, and invalid
forced zones return an error body with HTTP 400.

## Conventions

WGS84 constants are `a = 6378137 m` and `f = 1/298.257223563`. The central
meridian easting is 500000 m, the north equator northing is 0 m, and the
central scale is 0.9996. Southern-hemisphere results add 10000000 m to the
northing.

## Code Layout

```text
internal/utm       ellipsoid, zone, transverse Mercator series, validation
internal/grid      zone grid generation and point helpers
internal/verify    central easting, k0, southern-pair checks and tables
internal/server    HTTP handlers and JSON responses
internal/cli       subcommand parsing and terminal output
example/           offline scenario JSON files
```

## Build and Test

```bash
export GOTOOLCHAIN=local CGO_ENABLED=0
go build ./...
go test ./...
```

The Dockerfile builds the server binary and starts it on port 8080.
