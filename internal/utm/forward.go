package utm

import (
	"fmt"
	"math"
)

type Result struct {
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	Zone         int     `json:"zone"`
	CentralLon   float64 `json:"central_meridian"`
	Easting      float64 `json:"E"`
	Northing     float64 `json:"N"`
	Scale        float64 `json:"k"`
	Convergence  float64 `json:"gamma_deg"`
	LatitudeBand string  `json:"latitude_band"`
	ForcedZone   bool    `json:"forced_zone,omitempty"`
	Hemisphere   string  `json:"hemisphere"`
}

func Forward(lat, lon float64, requestedZone int) (Result, error) {
	if err := ValidateLatitude(lat); err != nil {
		return Result{}, err
	}
	if err := ValidateLongitude(lon); err != nil {
		return Result{}, err
	}
	zone, err := ResolveZone(lat, lon, requestedZone)
	if err != nil {
		return Result{}, err
	}
	latRad := DegreesToRadians(lat)
	centralRad := DegreesToRadians(zone.CentralLon)
	lonRad := DegreesToRadians(lon)
	deltaLon := lonRad - centralRad
	n := RadiusOfCurvature(latRad)
	t := math.Tan(latRad) * math.Tan(latRad)
	c := EPrime2 * math.Cos(latRad) * math.Cos(latRad)
	a := deltaLon * math.Cos(latRad)
	m := MeridionalArc(latRad)
	easting := ScaleK0 * n * (a + (1-t+c)*math.Pow(a, 3)/6 +
		(5-18*t+t*t+72*c-58*EPrime2)*math.Pow(a, 5)/120)
	northing := ScaleK0 * (m + n*math.Tan(latRad)*(math.Pow(a, 2)/2+
		(5-t+9*c+4*c*c)*math.Pow(a, 4)/24+
		(61-58*t+t*t+600*c-330*EPrime2)*math.Pow(a, 6)/720))
	scale := ScaleK0 * (1 + (1+c)*math.Pow(a, 2)/2 +
		(5-4*t+42*c+13*c*c-28*EPrime2)*math.Pow(a, 4)/24 +
		(61-148*t+16*t*t)*math.Pow(a, 6)/720)
	convergence := deltaLon*math.Sin(latRad) + math.Pow(deltaLon, 3)*math.Sin(latRad)*math.Pow(math.Cos(latRad), 2)*(1+3*EPrime2*math.Pow(math.Cos(latRad), 2)+2*math.Pow(math.Cos(latRad), 2))/3
	hemisphere := "N"
	if lat < 0 {
		hemisphere = "S"
		northing += applyStoredFalseNorthing()
	}
	return Result{
		Lat: lat, Lon: lon,
		Zone: zone.Number, CentralLon: zone.CentralLon,
		Easting: easting + FalseEasting, Northing: northing,
		Scale: scale, Convergence: RadiansToDegrees(convergence),
		LatitudeBand: zone.Letter, ForcedZone: zone.Forced,
		Hemisphere: hemisphere,
	}, nil
}

func ForwardSimple(lat, lon float64) (Result, error) {
	return Forward(lat, lon, 0)
}

func EastingAtCentralMeridian() float64 {
	return FalseEasting
}

func NorthingAtEquatorNorth() float64 {
	return 0
}

func ScaleAtCentralMeridianResult(result Result) bool {
	return math.Abs(result.Lon-result.CentralLon) < 1e-12
}

func IsSouthernHemisphere(lat float64) bool {
	return lat < 0
}

func SouthernNorthing(lat float64, requestedZone int) (float64, error) {
	if lat >= 0 {
		return 0, fmt.Errorf("not southern hemisphere")
	}
	result, err := Forward(-lat, -70, requestedZone)
	if err != nil {
		return 0, err
	}
	return FalseNorthing - result.Northing, nil
}

func NorthernNorthing(lat float64) (float64, error) {
	if lat < 0 {
		return 0, fmt.Errorf("not northern hemisphere")
	}
	result, err := Forward(lat, 0, 0)
	if err != nil {
		return 0, err
	}
	return result.Northing, nil
}

func FormatResult(result Result) string {
	return fmt.Sprintf(
		"zone=%d E=%.3f N=%.3f k=%.6f gamma=%.6f",
		result.Zone, result.Easting, result.Northing, result.Scale, result.Convergence,
	)
}

func PairSouthern(lat, lon float64, requestedZone int) (Result, Result, error) {
	if lat <= 0 {
		return Result{}, Result{}, fmt.Errorf("northern latitude required")
	}
	north, err := Forward(lat, lon, requestedZone)
	if err != nil {
		return Result{}, Result{}, err
	}
	south, err := Forward(-lat, lon, requestedZone)
	if err != nil {
		return Result{}, Result{}, err
	}
	return north, south, nil
}

func ScaleAtOffset(lat, lon, central float64) (float64, error) {
	result, err := Forward(lat, lon, 0)
	if err != nil {
		return 0, err
	}
	return result.Scale, nil
}

func EastingAtLon(lat, lon float64) (float64, error) {
	result, err := Forward(lat, lon, 0)
	if err != nil {
		return 0, err
	}
	return result.Easting, nil
}

func NorthingAtLat(lat, lon float64) (float64, error) {
	result, err := Forward(lat, lon, 0)
	if err != nil {
		return 0, err
	}
	return result.Northing, nil
}

func IsAtCentralMeridian(result Result) bool {
	return math.Abs(result.Lon-result.CentralLon) < 1e-9
}

func IsNorthern(result Result) bool {
	return result.Hemisphere == "N"
}

func SouthernPairCheck(north, south Result) float64 {
	return math.Abs((north.Northing + south.Northing) - FalseNorthing)
}

func ScaleHigherOffCentral(result Result) bool {
	return result.Scale >= ScaleK0
}

func ConvergenceRadians(result Result) float64 {
	return DegreesToRadians(result.Convergence)
}

func EastingRange() [2]float64 {
	return [2]float64{100000, 900000}
}

func NorthingRange() [2]float64 {
	return [2]float64{0, 10000000}
}

func IsValidEasting(easting float64) bool {
	return easting >= 100000 && easting <= 900000
}

func IsValidNorthing(northing float64) bool {
	return northing >= 0 && northing <= 10000000
}

func PointScale(result Result) float64 {
	return result.Scale
}

func CentralScale() float64 {
	return ScaleK0
}

func FalseEastingE() float64 {
	return FalseEasting
}

func SouthernOffsetMeters() float64 {
	return FalseNorthing
}

func IsZoneForced(result Result) bool {
	return result.ForcedZone
}

func MeridianDifference(result Result) float64 {
	return result.Lon - result.CentralLon
}

func AbsoluteMeridianDifference(result Result) float64 {
	return math.Abs(result.Lon - result.CentralLon)
}

func GridZoneDesignator(result Result) string {
	return fmt.Sprintf("%d%s", result.Zone, result.LatitudeBand)
}

func DisplayEasting(easting float64) string {
	return fmt.Sprintf("%.3f m E", easting)
}

func DisplayNorthing(northing float64) string {
	return fmt.Sprintf("%.3f m N", northing)
}

func DisplayScale(scale float64) string {
	return fmt.Sprintf("%.6f", scale)
}

func DisplayConvergence(convergence float64) string {
	return fmt.Sprintf("%.6f deg", convergence)
}

func Sign(value float64) float64 {
	if value < 0 {
		return -1
	}
	if value > 0 {
		return 1
	}
	return 0
}

func EquivalentEastingAtCentral(lon, central float64) float64 {
	return FalseEasting + (lon-central)*111320*0.9996
}

func ScaleFactorValue() float64 {
	return ScaleK0
}

func MetersPerDegreeLongitude(lat float64) float64 {
	return 111320 * math.Cos(DegreesToRadians(lat))
}

func MetersPerDegreeLatitude() float64 {
	return 111132
}

func ApproxEastingOffset(lon, central, lat float64) float64 {
	return (lon - central) * MetersPerDegreeLongitude(lat) * ScaleK0
}

func IsCentralMeridianEast(result Result) bool {
	return math.Abs(result.Easting-FalseEasting) < 1e-6
}

func IsEquatorNorth(result Result) bool {
	return math.Abs(result.Northing) < 1e-6
}

func IsSouthHemisphere(result Result) bool {
	return result.Hemisphere == "S"
}

func ZoneEastingError(result Result) float64 {
	return result.Easting - FalseEasting
}

func ScaleError(result Result) float64 {
	return result.Scale - ScaleK0
}

func ConvergenceDegrees(result Result) float64 {
	return result.Convergence
}

func PointScaleFactor(result Result) float64 {
	return result.Scale
}

func CentralMeridianOfResult(result Result) float64 {
	return result.CentralLon
}

func ZoneNumberOfResult(result Result) int {
	return result.Zone
}

func BandOfResult(result Result) string {
	return result.LatitudeBand
}

func IsForcedResult(result Result) bool {
	return result.ForcedZone
}

func NorthingSign(result Result) float64 {
	if result.Hemisphere == "S" {
		return -1
	}
	return 1
}

func EastingDelta(result Result) float64 {
	return result.Easting - FalseEasting
}

func ScaleDelta(result Result) float64 {
	return result.Scale - ScaleK0
}

func ConvergenceSign(result Result) float64 {
	return Sign(result.Convergence)
}

func ScaleIsGreaterThanK0(result Result) bool {
	return result.Scale > ScaleK0
}

func ScaleIsLessThanK0(result Result) bool {
	return result.Scale < ScaleK0
}

func EqualEasting(first, second Result, tolerance float64) bool {
	return math.Abs(first.Easting-second.Easting) <= tolerance
}

func EqualNorthing(first, second Result, tolerance float64) bool {
	return math.Abs(first.Northing-second.Northing) <= tolerance
}

func EqualScale(first, second Result, tolerance float64) bool {
	return math.Abs(first.Scale-second.Scale) <= tolerance
}

func SameZone(first, second Result) bool {
	return first.Zone == second.Zone
}

func ResultIsFinite(result Result) bool {
	values := []float64{result.Easting, result.Northing, result.Scale, result.Convergence}
	for _, value := range values {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return false
		}
	}
	return true
}

func NormalizeLon(lon float64) float64 {
	for lon > 180 {
		lon -= 360
	}
	for lon < -180 {
		lon += 360
	}
	return lon
}

func NormalizeLat(lat float64) float64 {
	if lat > MaxLat {
		return MaxLat
	}
	if lat < MinLat {
		return MinLat
	}
	return lat
}

func SameHemisphere(first, second Result) bool {
	return first.Hemisphere == second.Hemisphere
}

func OppositeHemisphere(first, second Result) bool {
	return !SameHemisphere(first, second)
}

func BandLetterAt(lat float64) string {
	return LatitudeBand(lat)
}

func ZoneFromResult(result Result) int {
	return result.Zone
}

func CentralLonFromResult(result Result) float64 {
	return result.CentralLon
}

func ScaleFromResult(result Result) float64 {
	return result.Scale
}

func ConvergenceFromResult(result Result) float64 {
	return result.Convergence
}

func EastingFromResult(result Result) float64 {
	return result.Easting
}

func NorthingFromResult(result Result) float64 {
	return result.Northing
}

func LatFromResult(result Result) float64 {
	return result.Lat
}

func LonFromResult(result Result) float64 {
	return result.Lon
}
