package utm

import "math"

const (
	A             = 6378137.0
	InvFlattening = 298.257223563
	ScaleK0       = 0.9996
	FalseEasting  = 500000.0
	FalseNorthing = 10000000.0
	MinLat        = -80.0
	MaxLat        = 84.0
	MinLon        = -180.0
	MaxLon        = 180.0
)

var (
	Flattening = 1 / InvFlattening
	E2         = Flattening * (2 - Flattening)
	EPrime2    = E2 / (1 - E2)
)

func MeridionalArc(latRad float64) float64 {
	e2 := E2
	e4 := e2 * e2
	e6 := e4 * e2
	return A * ((1-e2/4-3*e4/64-5*e6/256)*latRad -
		(3*e2/8+3*e4/32+45*e6/1024)*math.Sin(2*latRad) +
		(15*e4/256+45*e6/1024)*math.Sin(4*latRad) -
		(35*e6/3072)*math.Sin(6*latRad))
}

func RadiusOfCurvature(latRad float64) float64 {
	sinLat := math.Sin(latRad)
	return A / math.Sqrt(1-E2*sinLat*sinLat)
}

func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func RadiansToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

func EllipsoidSummary() string {
	return "WGS84: a=6378137 m, f=1/298.257223563"
}

func ScaleAtCentralMeridian() float64 {
	return ScaleK0
}

func CentralMeridianOffset() float64 {
	return 3.0
}

func ZoneWidth() float64 {
	return 6.0
}

func ZoneStartOffset() float64 {
	return 180.0
}

func MaxLatitude() float64 {
	return MaxLat
}

func MinLatitude() float64 {
	return MinLat
}

func MaxLongitude() float64 {
	return MaxLon
}

func MinLongitude() float64 {
	return MinLon
}

func SemimajorAxis() float64 {
	return A
}

func FlatteningValue() float64 {
	return Flattening
}

func EccentricitySquared() float64 {
	return E2
}

func SecondEccentricitySquared() float64 {
	return EPrime2
}

func ScaleFactor() float64 {
	return ScaleK0
}

func EquatorialRadius() float64 {
	return A
}

func FalseEastingValue() float64 {
	return FalseEasting
}

func SouthernOffset() float64 {
	return FalseNorthing
}

func ValidLatitudeRange() [2]float64 {
	return [2]float64{MinLat, MaxLat}
}

func ValidLongitudeRange() [2]float64 {
	return [2]float64{MinLon, MaxLon}
}
