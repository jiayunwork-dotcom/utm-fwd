package grid

import (
	"testing"
)

func TestGenerateGrid(t *testing.T) {
	grid, err := Generate(0, 10, 110, 130, 5, 5)
	if err != nil {
		t.Fatal(err)
	}
	if !IsGridNonEmpty(grid) || !GridOK(grid) {
		t.Fatalf("grid = %+v", grid)
	}
	if len(grid.Points) != grid.Rows*grid.Cols {
		t.Fatalf("rows=%d cols=%d points=%d", grid.Rows, grid.Cols, len(grid.Points))
	}
}

func TestFromLatLon(t *testing.T) {
	point, err := FromLatLon(39.9, 116.4)
	if err != nil {
		t.Fatal(err)
	}
	if point.Zone != 50 {
		t.Fatalf("zone = %d", point.Zone)
	}
	if point.Easting < 400000 || point.Easting > 600000 {
		t.Fatalf("E = %g", point.Easting)
	}
}

func TestZonesIn(t *testing.T) {
	grid, _ := Generate(-10, 10, -10, 10, 5, 5)
	zones := ZonesIn(grid.Points)
	if len(zones) == 0 {
		t.Fatal("no zones")
	}
}

func TestCentralMeridianPoints(t *testing.T) {
	grid, err := Generate(0, 10, 117, 118, 5, 1)
	if err != nil {
		t.Fatal(err)
	}
	points := CentralMeridianPoints(grid.Points)
	if len(points) == 0 {
		t.Fatalf("no central points; grid=%+v", grid)
	}
}

func TestEquatorPoints(t *testing.T) {
	grid, err := Generate(-1, 1, 110, 120, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(EquatorPoints(grid.Points)) == 0 {
		t.Fatal("no equator points")
	}
}

func TestSouthernNorthern(t *testing.T) {
	grid, _ := Generate(-10, 10, 110, 120, 5, 5)
	if len(SouthernPoints(grid.Points)) == 0 || len(NorthernPoints(grid.Points)) == 0 {
		t.Fatal("expected both hemispheres")
	}
}

func TestSummary(t *testing.T) {
	grid, _ := Generate(0, 10, 110, 120, 5, 5)
	if Summary(grid) == "" {
		t.Fatal("empty summary")
	}
	if Text(grid) == "" {
		t.Fatal("empty text")
	}
}

func TestScaleAtPoint(t *testing.T) {
	point, _ := FromLatLon(39.9, 116.4)
	scale, err := ScaleAtPoint(point)
	if err != nil {
		t.Fatal(err)
	}
	if scale <= 0.999 || scale > 1.001 {
		t.Fatalf("scale = %g", scale)
	}
}

func TestPointHelpers(t *testing.T) {
	point, _ := FromLatLon(39.9, 116.4)
	if !IsNorthern(point) || IsSouthern(point) {
		t.Fatal("hemisphere helper wrong")
	}
	if IsCentral(point) {
		t.Fatal("beijing incorrectly central")
	}
	if IsEquator(point) {
		t.Fatal("beijing incorrectly equator")
	}
}
