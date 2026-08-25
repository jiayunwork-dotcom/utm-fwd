package cli

import (
	"fmt"
	"os"

	"utm-fwd/internal/utm"
)

func Run(args []string) int {
	if len(args) == 0 {
		return runServe([]string{})
	}
	switch args[0] {
	case "forward":
		return runForward(args[1:])
	case "zone":
		return runZone(args[1:])
	case "example":
		return runExample(args[1:])
	case "serve":
		return runServe(args[1:])
	case "help", "-h", "--help":
		printHelp()
		return 0
	case "version":
		fmt.Println("utm-fwd 1.0.0")
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", args[0])
		printHelp()
		return 2
	}
}

func printHelp() {
	fmt.Println(`utm-fwd: WGS84 UTM forward projection calculator

Usage:
  utm-fwd                          start HTTP server on :8080
  utm-fwd forward -lat 39.9 -lon 116.4
  utm-fwd forward -lat 39.9 -lon 116.4 -zone 50
  utm-fwd zone -lon 116.4
  utm-fwd serve -addr :8080

HTTP:
  POST /api/forward  {"lat":39.9,"lon":116.4}
  POST /api/zone     {"lon":116.4}`)
}

func fail(err error) int {
	fmt.Fprintln(os.Stderr, "error:", err)
	return 1
}

func runForward(args []string) int {
	fs := flagSet("forward")
	lat := fs.Float64("lat", 0, "latitude")
	lon := fs.Float64("lon", 0, "longitude")
	zone := fs.Int("zone", 0, "requested zone (optional)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	result, err := utm.Forward(*lat, *lon, *zone)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}

func runZone(args []string) int {
	fs := flagSet("zone")
	lon := fs.Float64("lon", 0, "longitude")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	zone, err := utm.ZoneForLongitude(*lon)
	if err != nil {
		return fail(err)
	}
	return printJSON(zone)
}
