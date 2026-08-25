package utm

import (
	"math"
	"testing"
)

func TestCentralEasting500000(t *testing.T) {
	for _, lat := range []float64{0, 20, 39.9, -30} {
		result, err := Forward(lat, 117, 0)
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(result.Easting-500000) > 1e-6 {
			t.Fatalf("lat=%g E=%g, want 500000", lat, result.Easting)
		}
	}
}

func TestEquatorNorthingZeroNorth(t *testing.T) {
	result, err := Forward(0, 117, 0)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(result.Northing) > 1e-6 {
		t.Fatalf("N=%g, want 0", result.Northing)
	}
}

func TestCentralScaleK0(t *testing.T) {
	result, err := Forward(39.9, 117, 0)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(result.Scale-0.9996) > 1e-9 {
		t.Fatalf("k=%g, want 0.9996", result.Scale)
	}
}

func TestZoneFormulaBeijing(t *testing.T) {
	zone, err := ZoneForLongitude(116.4)
	if err != nil {
		t.Fatal(err)
	}
	if zone.Number != 50 {
		t.Fatalf("zone=%d, want 50", zone.Number)
	}
	if math.Abs(zone.CentralLon-117) > 1e-9 {
		t.Fatalf("central=%g, want 117", zone.CentralLon)
	}
}

func TestSouthernPair(t *testing.T) {
	north, south, err := PairSouthern(39.9, 116.4, 0)
	if err != nil {
		t.Fatal(err)
	}
	sum := north.Northing + south.Northing
	if math.Abs(sum-10000000) > 1e-6 {
		t.Fatalf("sum=%g, want 10000000", sum)
	}
	if math.Abs(north.Easting-south.Easting) > 1e-6 {
		t.Fatalf("eastings differ: %g vs %g", north.Easting, south.Easting)
	}
}

func TestScaleIncreasesOffCentral(t *testing.T) {
	central, err := Forward(39.9, 117, 0)
	if err != nil {
		t.Fatal(err)
	}
	offset, err := Forward(39.9, 116.4, 0)
	if err != nil {
		t.Fatal(err)
	}
	if offset.Scale <= central.Scale {
		t.Fatalf("offset scale %g not above central %g", offset.Scale, central.Scale)
	}
}

func TestValidationRejectsBadLatLon(t *testing.T) {
	if err := ValidateLatitude(85); err == nil {
		t.Fatal("accepted lat 85")
	}
	if err := ValidateLatitude(-81); err == nil {
		t.Fatal("accepted lat -81")
	}
	if err := ValidateLongitude(181); err == nil {
		t.Fatal("accepted lon 181")
	}
	if err := ValidateLongitude(-181); err == nil {
		t.Fatal("accepted lon -181")
	}
}

func TestForcedZoneFlag(t *testing.T) {
	result, err := Forward(39.9, 116.4, 51)
	if err != nil {
		t.Fatal(err)
	}
	if !result.ForcedZone {
		t.Fatal("forced zone not flagged")
	}
	auto, err := Forward(39.9, 116.4, 0)
	if err != nil {
		t.Fatal(err)
	}
	if auto.ForcedZone {
		t.Fatal("auto zone incorrectly flagged forced")
	}
}

func TestLatitudeBand(t *testing.T) {
	if LatitudeBand(39.9) != "S" {
		t.Fatalf("band = %q", LatitudeBand(39.9))
	}
	if LatitudeBand(45) != "T" {
		t.Fatalf("band = %q", LatitudeBand(45))
	}
	if LatitudeBand(-10) != "L" {
		t.Fatalf("band = %q", LatitudeBand(-10))
	}
}

func TestBeijingEastingRange(t *testing.T) {
	result, err := Forward(39.9, 116.4, 0)
	if err != nil {
		t.Fatal(err)
	}
	if result.Easting < 400000 || result.Easting > 600000 {
		t.Fatalf("E=%g outside 4e5..6e5", result.Easting)
	}
}

func TestConvergenceFinite(t *testing.T) {
	result, err := Forward(39.9, 116.4, 0)
	if err != nil {
		t.Fatal(err)
	}
	if math.IsNaN(result.Convergence) || math.IsInf(result.Convergence, 0) {
		t.Fatalf("gamma=%g", result.Convergence)
	}
}

func TestMeridionalArcPositive(t *testing.T) {
	arc := MeridionalArc(DegreesToRadians(39.9))
	if arc <= 0 {
		t.Fatalf("arc=%g", arc)
	}
}

func TestConstants(t *testing.T) {
	if A != 6378137 {
		t.Fatalf("a=%g", A)
	}
	if ScaleK0 != 0.9996 {
		t.Fatalf("k0=%g", ScaleK0)
	}
}
