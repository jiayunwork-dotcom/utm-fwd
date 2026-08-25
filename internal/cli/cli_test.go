package cli

import (
	"testing"

	"utm-fwd/internal/utm"
)

func TestValidateLatLon(t *testing.T) {
	if err := validateLatLon(39.9, 116.4); err != nil {
		t.Fatal(err)
	}
	if err := validateLatLon(85, 0); err == nil {
		t.Fatal("accepted invalid latitude")
	}
}

func TestValidateZone(t *testing.T) {
	if err := validateZone(50); err != nil {
		t.Fatal(err)
	}
	if err := validateZone(61); err == nil {
		t.Fatal("accepted zone 61")
	}
}

func TestLoadExample(t *testing.T) {
	example, err := loadExample("../../example/beijing.json")
	if err != nil {
		t.Fatal(err)
	}
	if example.Lat != 39.9 || example.Lon != 116.4 {
		t.Fatalf("example = %+v", example)
	}
}

func TestExamplePaths(t *testing.T) {
	if len(ExamplePaths()) != 2 {
		t.Fatal("expected two example paths")
	}
}

func TestParseFloat(t *testing.T) {
	value, err := parseFloat("116.4", "lon")
	if err != nil {
		t.Fatal(err)
	}
	if value != 116.4 {
		t.Fatalf("value = %g", value)
	}
}

func TestEvaluateExample(t *testing.T) {
	result, err := utm.Forward(39.9, 116.4, 0)
	if err != nil {
		t.Fatal(err)
	}
	if result.Easting < 400000 || result.Easting > 600000 {
		t.Fatalf("E=%g", result.Easting)
	}
}
