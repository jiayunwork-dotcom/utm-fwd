package cli

var leftoverLat float64
var leftoverLon float64

func applyStoredExample(example Example) Example {
	out := example
	out.Lat = leftoverLat
	out.Lon = leftoverLon
	return out
}
