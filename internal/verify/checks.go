package verify

import (
	"fmt"
	"math"

	"utm-fwd/internal/utm"
)

type Check struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func CheckCentralEasting(lat, lon float64) Check {
	zone, err := utm.ZoneForLongitude(lon)
	if err != nil {
		return Check{Name: "central-easting", OK: false, Message: err.Error()}
	}
	result, err := utm.Forward(lat, zone.CentralLon, 0)
	if err != nil {
		return Check{Name: "central-easting", OK: false, Message: err.Error()}
	}
	ok := math.Abs(result.Easting-utm.FalseEastingValue()) < 1e-6
	return Check{
		Name:    "central-easting",
		OK:      ok,
		Message: fmt.Sprintf("E=%.6f", result.Easting),
	}
}

func CheckEquatorNorthing(lon float64) Check {
	result, err := utm.Forward(0, lon, 0)
	if err != nil {
		return Check{Name: "equator-northing", OK: false, Message: err.Error()}
	}
	ok := math.Abs(result.Northing) < 1e-6
	return Check{
		Name:    "equator-northing",
		OK:      ok,
		Message: fmt.Sprintf("N=%.6f", result.Northing),
	}
}

func CheckCentralScale(lat, lon float64) Check {
	zone, err := utm.ZoneForLongitude(lon)
	if err != nil {
		return Check{Name: "central-scale", OK: false, Message: err.Error()}
	}
	result, err := utm.Forward(lat, zone.CentralLon, 0)
	if err != nil {
		return Check{Name: "central-scale", OK: false, Message: err.Error()}
	}
	ok := math.Abs(result.Scale-utm.ScaleK0) < 1e-9
	return Check{
		Name:    "central-scale",
		OK:      ok,
		Message: fmt.Sprintf("k=%.9f", result.Scale),
	}
}

func CheckScaleIncreases(lat, lon float64) Check {
	zone, err := utm.ZoneForLongitude(lon)
	if err != nil {
		return Check{Name: "scale-increase", OK: false, Message: err.Error()}
	}
	central, _ := utm.Forward(lat, zone.CentralLon, 0)
	offset, _ := utm.Forward(lat, lon, 0)
	ok := offset.Scale > central.Scale
	return Check{
		Name:    "scale-increase",
		OK:      ok,
		Message: fmt.Sprintf("central=%.9f offset=%.9f", central.Scale, offset.Scale),
	}
}

func CheckSouthernPair(lat, lon float64) Check {
	if lat <= 0 {
		return Check{Name: "southern-pair", OK: false, Message: "need northern latitude"}
	}
	north, south, err := utm.PairSouthern(lat, lon, 0)
	if err != nil {
		return Check{Name: "southern-pair", OK: false, Message: err.Error()}
	}
	sum := north.Northing + south.Northing
	ok := math.Abs(sum-utm.FalseNorthing) < 1e-6
	return Check{
		Name:    "southern-pair",
		OK:      ok,
		Message: fmt.Sprintf("north=%.6f south=%.6f sum=%.6f", north.Northing, south.Northing, sum),
	}
}

func CheckZoneFormula(lon float64, expected int) Check {
	zone, err := utm.ZoneForLongitude(lon)
	if err != nil {
		return Check{Name: "zone-formula", OK: false, Message: err.Error()}
	}
	ok := zone.Number == expected
	return Check{
		Name:    "zone-formula",
		OK:      ok,
		Message: fmt.Sprintf("zone=%d expected=%d", zone.Number, expected),
	}
}

func CheckCentralMeridian(number int, expected float64) Check {
	central := utm.CentralMeridian(number)
	ok := math.Abs(central-expected) < 1e-9
	return Check{
		Name:    "central-meridian",
		OK:      ok,
		Message: fmt.Sprintf("central=%g expected=%g", central, expected),
	}
}

func CheckForcedZone(lat, lon float64, requested int) Check {
	result, err := utm.Forward(lat, lon, requested)
	if err != nil {
		return Check{Name: "forced-zone", OK: false, Message: err.Error()}
	}
	expected := requested != utm.ZoneNumberFromLongitude(lon)
	ok := result.ForcedZone == expected
	return Check{
		Name:    "forced-zone",
		OK:      ok,
		Message: fmt.Sprintf("forced=%v expected=%v", result.ForcedZone, expected),
	}
}

func CheckValidation(lat, lon float64) Check {
	err := utm.ValidateLatLon(lat, lon)
	ok := err != nil
	return Check{Name: "validation", OK: ok, Message: fmt.Sprintf("err=%v", err)}
}

func RunAll() []Check {
	return []Check{
		CheckCentralEasting(39.9, 116.4),
		CheckEquatorNorthing(117),
		CheckCentralScale(39.9, 116.4),
		CheckScaleIncreases(39.9, 116.4),
		CheckSouthernPair(39.9, 116.4),
		CheckZoneFormula(116.4, 50),
		CheckCentralMeridian(50, 117),
		CheckForcedZone(39.9, 116.4, 51),
		CheckValidation(85, 0),
	}
}

func AllPass(checks []Check) bool {
	for _, check := range checks {
		if !check.OK {
			return false
		}
	}
	return true
}

func FormatChecks(checks []Check) string {
	out := ""
	for _, check := range checks {
		state := "PASS"
		if !check.OK {
			state = "FAIL"
		}
		out += fmt.Sprintf("%-18s %s %s\n", check.Name, state, check.Message)
	}
	return out
}

func CheckZoneOffByOne(lon float64) Check {
	auto, err := utm.ZoneForLongitude(lon)
	if err != nil {
		return Check{Name: "zone-off-by-one", OK: false, Message: err.Error()}
	}
	wrong := auto.Number - 1
	if wrong < 1 {
		wrong = auto.Number + 1
	}
	result, err := utm.Forward(0, lon, wrong)
	if err != nil {
		return Check{Name: "zone-off-by-one", OK: false, Message: err.Error()}
	}
	ok := math.Abs(result.Easting-utm.FalseEastingValue()) > 100000
	return Check{
		Name:    "zone-off-by-one",
		OK:      ok,
		Message: fmt.Sprintf("wrong zone %d E=%.3f", wrong, result.Easting),
	}
}

func CheckBand(lat float64, expected string) Check {
	band := utm.LatitudeBand(lat)
	ok := band == expected
	return Check{Name: "band", OK: ok, Message: fmt.Sprintf("band=%s expected=%s", band, expected)}
}

func CheckScaleRange(lat, lon float64) Check {
	result, err := utm.Forward(lat, lon, 0)
	if err != nil {
		return Check{Name: "scale-range", OK: false, Message: err.Error()}
	}
	ok := result.Scale >= 0.999 && result.Scale <= 1.001
	return Check{Name: "scale-range", OK: ok, Message: fmt.Sprintf("k=%.6f", result.Scale)}
}

func CheckConvergenceFinite(lat, lon float64) Check {
	result, err := utm.Forward(lat, lon, 0)
	if err != nil {
		return Check{Name: "convergence-finite", OK: false, Message: err.Error()}
	}
	ok := !math.IsNaN(result.Convergence) && !math.IsInf(result.Convergence, 0)
	return Check{Name: "convergence-finite", OK: ok, Message: fmt.Sprintf("gamma=%.6f", result.Convergence)}
}

func CheckEastingPlausible(lat, lon float64) Check {
	result, err := utm.Forward(lat, lon, 0)
	if err != nil {
		return Check{Name: "easting-plausible", OK: false, Message: err.Error()}
	}
	ok := result.Easting >= 100000 && result.Easting <= 900000
	return Check{Name: "easting-plausible", OK: ok, Message: fmt.Sprintf("E=%.3f", result.Easting)}
}

func CheckNorthingPlausible(lat, lon float64) Check {
	result, err := utm.Forward(lat, lon, 0)
	if err != nil {
		return Check{Name: "northing-plausible", OK: false, Message: err.Error()}
	}
	ok := result.Northing >= 0 && result.Northing <= 10000000
	return Check{Name: "northing-plausible", OK: ok, Message: fmt.Sprintf("N=%.3f", result.Northing)}
}

func CheckForcedFlag(result utm.Result) Check {
	ok := result.ForcedZone
	return Check{Name: "forced-flag", OK: ok, Message: fmt.Sprintf("forced=%v", result.ForcedZone)}
}

func CheckAutoFlag(result utm.Result) Check {
	ok := !result.ForcedZone
	return Check{Name: "auto-flag", OK: ok, Message: fmt.Sprintf("forced=%v", result.ForcedZone)}
}

func CheckSameScalePair(first, second utm.Result, tolerance float64) Check {
	ok := math.Abs(first.Scale-second.Scale) <= tolerance
	return Check{Name: "scale-pair", OK: ok, Message: fmt.Sprintf("%.9f vs %.9f", first.Scale, second.Scale)}
}

func CheckSameZonePair(first, second utm.Result) Check {
	ok := first.Zone == second.Zone
	return Check{Name: "zone-pair", OK: ok, Message: fmt.Sprintf("%d vs %d", first.Zone, second.Zone)}
}

func CheckOppositeHemispherePair(first, second utm.Result) Check {
	ok := first.Hemisphere != second.Hemisphere
	return Check{Name: "hemisphere-pair", OK: ok, Message: fmt.Sprintf("%s vs %s", first.Hemisphere, second.Hemisphere)}
}

func CheckBandPair(first, second utm.Result) Check {
	ok := first.LatitudeBand == second.LatitudeBand
	return Check{Name: "band-pair", OK: ok, Message: fmt.Sprintf("%s vs %s", first.LatitudeBand, second.LatitudeBand)}
}

func CheckEastingPair(first, second utm.Result, tolerance float64) Check {
	ok := math.Abs(first.Easting-second.Easting) <= tolerance
	return Check{Name: "easting-pair", OK: ok, Message: fmt.Sprintf("%.6f vs %.6f", first.Easting, second.Easting)}
}

func CheckNorthingPairSum(first, second utm.Result, expected float64) Check {
	sum := first.Northing + second.Northing
	ok := math.Abs(sum-expected) < 1e-6
	return Check{Name: "northing-pair-sum", OK: ok, Message: fmt.Sprintf("sum=%.6f expected=%.6f", sum, expected)}
}

func CheckHemisphere(first, second utm.Result) Check {
	if first.Hemisphere != second.Hemisphere {
		return Check{Name: "hemisphere", OK: false, Message: "hemispheres differ"}
	}
	return Check{Name: "hemisphere", OK: true, Message: first.Hemisphere}
}

func CheckBandSame(first, second utm.Result) Check {
	if first.LatitudeBand != second.LatitudeBand {
		return Check{Name: "band-same", OK: false, Message: "bands differ"}
	}
	return Check{Name: "band-same", OK: true, Message: first.LatitudeBand}
}

func CheckZoneSame(first, second utm.Result) Check {
	if first.Zone != second.Zone {
		return Check{Name: "zone-same", OK: false, Message: "zones differ"}
	}
	return Check{Name: "zone-same", OK: true, Message: fmt.Sprintf("%d", first.Zone)}
}

func CheckEastingSame(first, second utm.Result, tolerance float64) Check {
	if math.Abs(first.Easting-second.Easting) > tolerance {
		return Check{Name: "easting-same", OK: false, Message: "eastings differ"}
	}
	return Check{Name: "easting-same", OK: true, Message: fmt.Sprintf("%.3f", first.Easting)}
}

func CheckNorthingSame(first, second utm.Result, tolerance float64) Check {
	if math.Abs(first.Northing-second.Northing) > tolerance {
		return Check{Name: "northing-same", OK: false, Message: "northings differ"}
	}
	return Check{Name: "northing-same", OK: true, Message: fmt.Sprintf("%.3f", first.Northing)}
}

func CheckScaleSame(first, second utm.Result, tolerance float64) Check {
	if math.Abs(first.Scale-second.Scale) > tolerance {
		return Check{Name: "scale-same", OK: false, Message: "scales differ"}
	}
	return Check{Name: "scale-same", OK: true, Message: fmt.Sprintf("%.9f", first.Scale)}
}

func CheckConvergenceSame(first, second utm.Result, tolerance float64) Check {
	if math.Abs(first.Convergence-second.Convergence) > tolerance {
		return Check{Name: "convergence-same", OK: false, Message: "convergences differ"}
	}
	return Check{Name: "convergence-same", OK: true, Message: fmt.Sprintf("%.6f", first.Convergence)}
}

func CheckZoneAutoEquator(lon float64) Check {
	zone, err := utm.ZoneForLongitude(lon)
	if err != nil {
		return Check{Name: "zone-auto-equator", OK: false, Message: err.Error()}
	}
	ok := zone.Number >= 1 && zone.Number <= 60
	return Check{Name: "zone-auto-equator", OK: ok, Message: fmt.Sprintf("zone=%d", zone.Number)}
}

func CheckScaleK0Exact(result utm.Result) Check {
	ok := math.Abs(result.Scale-utm.ScaleK0) < 1e-9
	return Check{Name: "scale-k0-exact", OK: ok, Message: fmt.Sprintf("k=%.12f", result.Scale)}
}

func CheckEasting500kExact(result utm.Result) Check {
	ok := math.Abs(result.Easting-500000) < 1e-6
	return Check{Name: "easting-500k", OK: ok, Message: fmt.Sprintf("E=%.6f", result.Easting)}
}

func CheckNorthingZeroExact(result utm.Result) Check {
	ok := math.Abs(result.Northing) < 1e-6
	return Check{Name: "northing-zero", OK: ok, Message: fmt.Sprintf("N=%.6f", result.Northing)}
}

func CheckSouthernAddsOffset(result utm.Result) Check {
	ok := result.Northing > 9000000
	return Check{Name: "southern-offset", OK: ok, Message: fmt.Sprintf("N=%.3f", result.Northing)}
}

func CheckEastingInRange(result utm.Result) Check {
	ok := result.Easting >= 100000 && result.Easting <= 900000
	return Check{Name: "easting-range", OK: ok, Message: fmt.Sprintf("E=%.3f", result.Easting)}
}

func CheckNorthingInRange(result utm.Result) Check {
	ok := result.Northing >= 0 && result.Northing <= 10000000
	return Check{Name: "northing-range", OK: ok, Message: fmt.Sprintf("N=%.3f", result.Northing)}
}

func CheckScaleInRange(result utm.Result) Check {
	ok := result.Scale >= 0.999 && result.Scale <= 1.001
	return Check{Name: "scale-range", OK: ok, Message: fmt.Sprintf("k=%.6f", result.Scale)}
}

func CheckConvergenceInRange(result utm.Result) Check {
	ok := math.Abs(result.Convergence) <= 10
	return Check{Name: "convergence-range", OK: ok, Message: fmt.Sprintf("gamma=%.6f", result.Convergence)}
}

func CheckZoneInRange(result utm.Result) Check {
	ok := result.Zone >= 1 && result.Zone <= 60
	return Check{Name: "zone-range", OK: ok, Message: fmt.Sprintf("zone=%d", result.Zone)}
}

func CheckBandNonEmpty(result utm.Result) Check {
	ok := result.LatitudeBand != ""
	return Check{Name: "band-nonempty", OK: ok, Message: result.LatitudeBand}
}

func CheckForcedZoneResult(result utm.Result) Check {
	return Check{Name: "forced-zone-result", OK: result.ForcedZone, Message: fmt.Sprintf("forced=%v", result.ForcedZone)}
}

func CheckAutoZoneResult(result utm.Result) Check {
	return Check{Name: "auto-zone-result", OK: !result.ForcedZone, Message: fmt.Sprintf("forced=%v", result.ForcedZone)}
}

func CheckFiniteResult(result utm.Result) Check {
	return Check{Name: "finite-result", OK: utm.ResultIsFinite(result), Message: fmt.Sprintf("%+v", result)}
}

func CheckValidResult(result utm.Result) Check {
	return Check{Name: "valid-result", OK: utm.ValidateResult(result) == nil, Message: fmt.Sprintf("%+v", result)}
}

func CheckZoneAutoResult(result utm.Result) Check {
	return Check{Name: "zone-auto-result", OK: result.Zone == utm.ZoneNumberFromLongitude(result.Lon), Message: fmt.Sprintf("%d", result.Zone)}
}

func CheckScaleAtCentralResult(result utm.Result) Check {
	return Check{Name: "scale-at-central", OK: math.Abs(result.Scale-utm.ScaleK0) < 1e-9, Message: fmt.Sprintf("k=%.12f", result.Scale)}
}

func CheckEastingAtCentralResult(result utm.Result) Check {
	return Check{Name: "easting-at-central", OK: math.Abs(result.Easting-500000) < 1e-6, Message: fmt.Sprintf("E=%.6f", result.Easting)}
}

func CheckNorthingAtEquatorResult(result utm.Result) Check {
	return Check{Name: "northing-at-equator", OK: math.Abs(result.Northing) < 1e-6, Message: fmt.Sprintf("N=%.6f", result.Northing)}
}
