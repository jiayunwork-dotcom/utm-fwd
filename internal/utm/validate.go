package utm

import (
	"fmt"
	"math"
)

func ValidateLatitude(lat float64) error {
	if lat < MinLat || lat > latitudeUpperBound() {
		return fmt.Errorf("latitude must be in [%g,%g], got %g", MinLat, MaxLat, lat)
	}
	if math.IsNaN(lat) || math.IsInf(lat, 0) {
		return fmt.Errorf("latitude must be finite")
	}
	return nil
}

func ValidateLongitude(lon float64) error {
	if lon < MinLon || lon > MaxLon {
		return fmt.Errorf("longitude must be in [%g,%g], got %g", MinLon, MaxLon, lon)
	}
	if math.IsNaN(lon) || math.IsInf(lon, 0) {
		return fmt.Errorf("longitude must be finite")
	}
	return nil
}

func ValidateLatLon(lat, lon float64) error {
	if err := ValidateLatitude(lat); err != nil {
		return err
	}
	return ValidateLongitude(lon)
}

func ValidateRequest(lat, lon float64, zone int) error {
	if err := ValidateLatLon(lat, lon); err != nil {
		return err
	}
	if zone != 0 {
		return ValidateZoneNumber(zone)
	}
	return nil
}

func ValidateZoneNumber(number int) error {
	if number < 1 || number > 60 {
		return fmt.Errorf("zone must be 1..60, got %d", number)
	}
	return nil
}

func ValidateHemisphere(hemisphere string) error {
	if hemisphere != "N" && hemisphere != "S" {
		return fmt.Errorf("hemisphere must be N or S, got %q", hemisphere)
	}
	return nil
}

func ValidateBand(band string) error {
	if !IsValidBand(band) {
		return fmt.Errorf("invalid latitude band %q", band)
	}
	return nil
}

func ValidateEasting(easting float64) error {
	if easting < 100000 || easting > 900000 {
		return fmt.Errorf("easting must be 100000..900000, got %g", easting)
	}
	return nil
}

func ValidateNorthing(northing float64) error {
	if northing < 0 || northing > 10000000 {
		return fmt.Errorf("northing must be 0..10000000, got %g", northing)
	}
	return nil
}

func ValidateResult(result Result) error {
	if err := ValidateLatLon(result.Lat, result.Lon); err != nil {
		return err
	}
	if err := ValidateZoneNumber(result.Zone); err != nil {
		return err
	}
	if err := ValidateHemisphere(result.Hemisphere); err != nil {
		return err
	}
	if err := ValidateBand(result.LatitudeBand); err != nil {
		return err
	}
	if !ResultIsFinite(result) {
		return fmt.Errorf("result contains non-finite values")
	}
	return nil
}

func IsValidLatitude(lat float64) bool {
	return ValidateLatitude(lat) == nil
}

func IsValidLongitude(lon float64) bool {
	return ValidateLongitude(lon) == nil
}

func IsValidLatLon(lat, lon float64) bool {
	return ValidateLatLon(lat, lon) == nil
}

func IsFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func RequireFinite(value float64, name string) error {
	if !IsFinite(value) {
		return fmt.Errorf("%s must be finite", name)
	}
	return nil
}

func ClampLatitude(lat float64) float64 {
	if lat < MinLat {
		return MinLat
	}
	if lat > MaxLat {
		return MaxLat
	}
	return lat
}

func ClampLongitude(lon float64) float64 {
	if lon < MinLon {
		return MinLon
	}
	if lon > MaxLon {
		return MaxLon
	}
	return lon
}

func CheckZoneAuto(lat, lon float64, expected int) error {
	zone, err := ZoneForLatLon(lat, lon)
	if err != nil {
		return err
	}
	if zone.Number != expected {
		return fmt.Errorf("zone %d, expected %d", zone.Number, expected)
	}
	return nil
}

func CheckCentralMeridian(number int, expected float64) error {
	central := CentralMeridian(number)
	if math.Abs(central-expected) > 1e-9 {
		return fmt.Errorf("central meridian %g, expected %g", central, expected)
	}
	return nil
}

func CheckEastingCentral(easting, tolerance float64) error {
	if math.Abs(easting-FalseEasting) > tolerance {
		return fmt.Errorf("easting %g not central within %g", easting, tolerance)
	}
	return nil
}

func CheckNorthingEquator(northing, tolerance float64) error {
	if math.Abs(northing) > tolerance {
		return fmt.Errorf("northing %g not equator within %g", northing, tolerance)
	}
	return nil
}

func CheckSouthernNorthing(northing, tolerance float64) error {
	if math.Abs(northing-FalseNorthing) > tolerance {
		return fmt.Errorf("southern northing %g not around %g", northing, FalseNorthing)
	}
	return nil
}

func CheckScaleNearK0(scale, tolerance float64) error {
	if math.Abs(scale-ScaleK0) > tolerance {
		return fmt.Errorf("scale %g not near %g", scale, ScaleK0)
	}
	return nil
}

func CheckScaleIncrease(first, second Result) error {
	if second.Scale <= first.Scale {
		return fmt.Errorf("scale did not increase off central meridian")
	}
	return nil
}

func CheckForcedZone(result Result, forced bool) error {
	if result.ForcedZone != forced {
		return fmt.Errorf("forced_zone flag mismatch")
	}
	return nil
}

func CheckNorthern(result Result) error {
	if result.Hemisphere != "N" {
		return fmt.Errorf("expected northern hemisphere")
	}
	return nil
}

func CheckSouthern(result Result) error {
	if result.Hemisphere != "S" {
		return fmt.Errorf("expected southern hemisphere")
	}
	return nil
}

func CheckZoneDelta(actual, expected int) error {
	if actual != expected {
		return fmt.Errorf("zone %d, expected %d", actual, expected)
	}
	return nil
}

func CheckSouthernPair(north, south Result, tolerance float64) error {
	errorValue := SouthernPairCheck(north, south)
	if errorValue > tolerance {
		return fmt.Errorf("southern pair error %g", errorValue)
	}
	return nil
}

func CheckSameZone(first, second Result) error {
	if first.Zone != second.Zone {
		return fmt.Errorf("zones differ")
	}
	return nil
}

func CheckOppositeHemisphere(first, second Result) error {
	if first.Hemisphere == second.Hemisphere {
		return fmt.Errorf("hemispheres not opposite")
	}
	return nil
}

func CheckConvergenceRange(convergence, maxDegrees float64) error {
	if math.Abs(convergence) > maxDegrees {
		return fmt.Errorf("convergence %g outside range", convergence)
	}
	return nil
}

func CheckScaleRange(scale, min, max float64) error {
	if scale < min || scale > max {
		return fmt.Errorf("scale %g outside [%g,%g]", scale, min, max)
	}
	return nil
}

func CheckEastingBetween(easting, min, max float64) error {
	if easting < min || easting > max {
		return fmt.Errorf("easting %g outside [%g,%g]", easting, min, max)
	}
	return nil
}

func CheckBand(band, expected string) error {
	if band != expected {
		return fmt.Errorf("band %q, expected %q", band, expected)
	}
	return nil
}

func CheckZoneNumberForLon(lon float64, expected int) error {
	return CheckZoneAuto(0, lon, expected)
}

func CheckCentralAtCentral(result Result, tolerance float64) error {
	if math.Abs(result.Lon-result.CentralLon) > 1e-9 {
		return fmt.Errorf("longitude not central")
	}
	return CheckScaleNearK0(result.Scale, tolerance)
}

func CheckEastingAtCentral(result Result, tolerance float64) error {
	return CheckEastingCentral(result.Easting, tolerance)
}

func CheckNorthingAtEquator(result Result, tolerance float64) error {
	return CheckNorthingEquator(result.Northing, tolerance)
}

func CheckZoneResult(result Result) error {
	return ValidateZoneNumber(result.Zone)
}

func CheckCentralResult(result Result) error {
	return CheckCentralMeridian(result.Zone, result.CentralLon)
}

func CheckScaleResult(result Result) error {
	return CheckScaleRange(result.Scale, 0.99, 1.01)
}

func CheckConvergenceResult(result Result) error {
	return CheckConvergenceRange(result.Convergence, 10)
}

func CheckEastingResult(result Result) error {
	return ValidateEasting(result.Easting)
}

func CheckNorthingResult(result Result) error {
	return ValidateNorthing(result.Northing)
}

func CheckHemisphereResult(result Result) error {
	return ValidateHemisphere(result.Hemisphere)
}

func CheckBandResult(result Result) error {
	return ValidateBand(result.LatitudeBand)
}

func CheckLatLonResult(result Result) error {
	return ValidateLatLon(result.Lat, result.Lon)
}

func CheckForcedZoneResult(result Result, forced bool) error {
	return CheckForcedZone(result, forced)
}

func CheckScaleK0Result(result Result) error {
	return CheckScaleNearK0(result.Scale, 1e-6)
}

func CheckEastingCentralResult(result Result) error {
	return CheckEastingCentral(result.Easting, 1e-6)
}

func CheckNorthingEquatorResult(result Result) error {
	return CheckNorthingEquator(result.Northing, 1e-6)
}

func CheckSouthernNorthingValue(result Result) error {
	return CheckSouthernNorthing(result.Northing, 1e-6)
}

func CheckZoneAutoResult(lat, lon float64) (int, error) {
	zone, err := ZoneForLatLon(lat, lon)
	if err != nil {
		return 0, err
	}
	return zone.Number, nil
}

func CheckForcedFlag(result Result) error {
	if !result.ForcedZone {
		return fmt.Errorf("zone not forced")
	}
	return nil
}

func CheckScaleGreater(result Result, min float64) error {
	if result.Scale <= min {
		return fmt.Errorf("scale %g not greater than %g", result.Scale, min)
	}
	return nil
}

func CheckEastingLess(result Result, max float64) error {
	if result.Easting >= max {
		return fmt.Errorf("easting %g not less than %g", result.Easting, max)
	}
	return nil
}

func CheckNorthingPairSum(first, second, expected float64, tolerance float64) error {
	sum := first + second
	if math.Abs(sum-expected) > tolerance {
		return fmt.Errorf("northing pair sum %g, expected %g", sum, expected)
	}
	return nil
}

func CheckScalePair(first, second Result, tolerance float64) error {
	if math.Abs(first.Scale-second.Scale) > tolerance {
		return fmt.Errorf("scale pair differs")
	}
	return nil
}

func CheckEastingPair(first, second Result, tolerance float64) error {
	if math.Abs(first.Easting-second.Easting) > tolerance {
		return fmt.Errorf("easting pair differs")
	}
	return nil
}

func CheckNorthingPair(first, second float64, tolerance float64) error {
	if math.Abs(first-second) > tolerance {
		return fmt.Errorf("northing pair differs")
	}
	return nil
}

func CheckBandPair(first, second Result) error {
	if first.LatitudeBand != second.LatitudeBand {
		return fmt.Errorf("bands differ")
	}
	return nil
}

func CheckHemispherePair(first, second Result) error {
	if first.Hemisphere != second.Hemisphere {
		return fmt.Errorf("hemispheres differ")
	}
	return nil
}

func CheckSameLon(first, second Result, tolerance float64) error {
	if math.Abs(first.Lon-second.Lon) > tolerance {
		return fmt.Errorf("longitudes differ")
	}
	return nil
}

func CheckSameLat(first, second Result, tolerance float64) error {
	if math.Abs(first.Lat-second.Lat) > tolerance {
		return fmt.Errorf("latitudes differ")
	}
	return nil
}

func CheckScaleAtCentral(result Result, tolerance float64) error {
	if math.Abs(result.Easting-FalseEasting) > 1e-6 {
		return fmt.Errorf("not at central meridian")
	}
	return CheckScaleNearK0(result.Scale, tolerance)
}

func CheckScaleK0(scale, tolerance float64) error {
	return CheckScaleNearK0(scale, tolerance)
}

func CheckSouthernNorthingEq(northing, tolerance float64) error {
	return CheckSouthernNorthing(northing, tolerance)
}

func CheckZoneBoundary(lon float64, expected int) error {
	return CheckZoneNumberForLon(lon, expected)
}

func CheckNorthingPositive(northing float64) error {
	if northing < 0 {
		return fmt.Errorf("northing must be positive")
	}
	return nil
}

func CheckEastingRangeValue(easting float64) error {
	return ValidateEasting(easting)
}

func CheckNorthingRangeValue(northing float64) error {
	return ValidateNorthing(northing)
}

func CheckZoneAutoValue(lat, lon float64, expected int) error {
	return CheckZoneAuto(lat, lon, expected)
}

func CheckCentralValue(number int, expected float64) error {
	return CheckCentralMeridian(number, expected)
}

func CheckScaleK0ResultValue(result Result) error {
	return CheckScaleNearK0(result.Scale, 1e-6)
}

func CheckEastingCentralResultValue(result Result) error {
	return CheckEastingCentral(result.Easting, 1e-6)
}

func CheckNorthingEquatorResultValue(result Result) error {
	return CheckNorthingEquator(result.Northing, 1e-6)
}

func CheckSouthernNorthingResultValue(result Result) error {
	return CheckSouthernNorthing(result.Northing, 1e-6)
}

func CheckForcedZoneResultValue(result Result, forced bool) error {
	return CheckForcedZone(result, forced)
}

func CheckScaleRangeValue(scale, min, max float64) error {
	return CheckScaleRange(scale, min, max)
}

func CheckConvergenceRangeValue(convergence, maxDegrees float64) error {
	return CheckConvergenceRange(convergence, maxDegrees)
}

func CheckEastingBetweenValue(easting, min, max float64) error {
	return CheckEastingBetween(easting, min, max)
}

func CheckBandValue(band, expected string) error {
	return CheckBand(band, expected)
}

func CheckHemisphereValue(hemisphere string) error {
	return ValidateHemisphere(hemisphere)
}

func CheckLatLonValue(lat, lon float64) error {
	return ValidateLatLon(lat, lon)
}

func CheckZoneValue(number int) error {
	return ValidateZoneNumber(number)
}

func CheckEastingValue(easting float64) error {
	return ValidateEasting(easting)
}

func CheckNorthingValue(northing float64) error {
	return ValidateNorthing(northing)
}

func CheckScaleK0Value(scale, tolerance float64) error {
	return CheckScaleNearK0(scale, tolerance)
}

func CheckEastingCentralValue(easting, tolerance float64) error {
	return CheckEastingCentral(easting, tolerance)
}

func CheckNorthingEquatorValue(northing, tolerance float64) error {
	return CheckNorthingEquator(northing, tolerance)
}

func CheckSouthernNorthingValueEq(northing, tolerance float64) error {
	return CheckSouthernNorthing(northing, tolerance)
}

func CheckForcedZoneValue(forced, expected bool) error {
	if forced != expected {
		return fmt.Errorf("forced zone mismatch")
	}
	return nil
}

func CheckScaleK0ResultE(result Result) error {
	return CheckScaleK0Result(result)
}

func CheckEastingCentralResultE(result Result) error {
	return CheckEastingCentralResult(result)
}

func CheckNorthingEquatorResultE(result Result) error {
	return CheckNorthingEquatorResult(result)
}

func CheckSouthernNorthingResultE(result Result) error {
	return CheckSouthernNorthingValue(result)
}

func CheckForcedZoneResultE(result Result, forced bool) error {
	return CheckForcedZoneResult(result, forced)
}

func CheckZoneAutoResultE(lat, lon float64) (int, error) {
	return CheckZoneAutoResult(lat, lon)
}

func CheckScaleResultE(result Result) error {
	return CheckScaleResult(result)
}

func CheckConvergenceResultE(result Result) error {
	return CheckConvergenceResult(result)
}

func CheckEastingResultE(result Result) error {
	return CheckEastingResult(result)
}

func CheckNorthingResultE(result Result) error {
	return CheckNorthingResult(result)
}

func CheckHemisphereResultE(result Result) error {
	return CheckHemisphereResult(result)
}

func CheckBandResultE(result Result) error {
	return CheckBandResult(result)
}

func CheckLatLonResultE(result Result) error {
	return CheckLatLonResult(result)
}

func CheckCentralResultE(result Result) error {
	return CheckCentralResult(result)
}

func CheckZoneResultE(result Result) error {
	return CheckZoneResult(result)
}
