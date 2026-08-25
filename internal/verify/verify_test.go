package verify

import (
	"testing"

	"utm-fwd/internal/utm"
)

func TestRunAllChecks(t *testing.T) {
	checks := RunAll()
	if !AllPass(checks) {
		for _, check := range checks {
			if !check.OK {
				t.Logf("%s: %s", check.Name, check.Message)
			}
		}
		t.Fatal("not all checks pass")
	}
}

func TestCheckCentralEasting(t *testing.T) {
	if !CheckCentralEasting(39.9, 116.4).OK {
		t.Fatal("central easting check failed")
	}
}

func TestCheckEquatorNorthing(t *testing.T) {
	if !CheckEquatorNorthing(117).OK {
		t.Fatal("equator northing check failed")
	}
}

func TestCheckCentralScale(t *testing.T) {
	if !CheckCentralScale(39.9, 116.4).OK {
		t.Fatal("central scale check failed")
	}
}

func TestCheckScaleIncreases(t *testing.T) {
	if !CheckScaleIncreases(39.9, 116.4).OK {
		t.Fatal("scale increase check failed")
	}
}

func TestCheckSouthernPair(t *testing.T) {
	if !CheckSouthernPair(39.9, 116.4).OK {
		t.Fatal("southern pair check failed")
	}
}

func TestCheckZoneFormula(t *testing.T) {
	if !CheckZoneFormula(116.4, 50).OK {
		t.Fatal("zone formula check failed")
	}
}

func TestCheckForcedZone(t *testing.T) {
	if !CheckForcedZone(39.9, 116.4, 51).OK {
		t.Fatal("forced zone check failed")
	}
}

func TestCheckValidation(t *testing.T) {
	if !CheckValidation(85, 0).OK {
		t.Fatal("validation check failed")
	}
}

func TestBuildTable(t *testing.T) {
	table, err := BuildTable([][2]float64{{0, 117}, {39.9, 116.4}})
	if err != nil {
		t.Fatal(err)
	}
	if !HasCentralRow(table) || !HasEquatorRow(table) {
		t.Fatal("table missing central/equator rows")
	}
}

func TestTableText(t *testing.T) {
	result, _ := utm.Forward(39.9, 116.4, 0)
	text := TableText([]utm.Result{result})
	if text == "" {
		t.Fatal("empty table text")
	}
}

func TestMedianEasting(t *testing.T) {
	r1, _ := utm.Forward(0, 0, 0)
	r2, _ := utm.Forward(39.9, 116.4, 0)
	if MedianEasting([]utm.Result{r1, r2}) <= 0 {
		t.Fatal("median not positive")
	}
}

func TestCheckScaleK0Result(t *testing.T) {
	result, _ := utm.Forward(0, 117, 0)
	if !CheckScaleK0Exact(result).OK {
		t.Fatal("k0 check failed")
	}
}

func TestCheckEastingCentralResult(t *testing.T) {
	result, _ := utm.Forward(0, 117, 0)
	if !CheckEasting500kExact(result).OK {
		t.Fatal("easting central check failed")
	}
}

func TestCheckNorthingEquatorResult(t *testing.T) {
	result, _ := utm.Forward(0, 117, 0)
	if !CheckNorthingZeroExact(result).OK {
		t.Fatal("northing equator check failed")
	}
}

func TestCSVText(t *testing.T) {
	result, _ := utm.Forward(39.9, 116.4, 0)
	if CSVText([]utm.Result{result}) == "" {
		t.Fatal("empty csv")
	}
}
