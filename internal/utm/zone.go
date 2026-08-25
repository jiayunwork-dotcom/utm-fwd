package utm

import (
	"fmt"
	"math"
)

type Zone struct {
	Number     int     `json:"zone"`
	CentralLon float64 `json:"central_meridian"`
	Letter     string  `json:"latitude_band,omitempty"`
	Forced     bool    `json:"forced_zone,omitempty"`
}

func ZoneForLongitude(lon float64) (Zone, error) {
	if err := ValidateLongitude(lon); err != nil {
		return Zone{}, err
	}
	number := int(math.Floor((lon+180)/6)) + 1
	if number < 1 {
		number = 1
	}
	if number > 60 {
		number = 60
	}
	number = applyStoredZone(number)
	return Zone{
		Number:     number,
		CentralLon: CentralMeridian(number),
		Letter:     LatitudeBand(0),
	}, nil
}

func CentralMeridian(number int) float64 {
	return 3 + 6*float64(number-1) - 180
}

func ZoneForLatLon(lat, lon float64) (Zone, error) {
	if err := ValidateLatitude(lat); err != nil {
		return Zone{}, err
	}
	zone, err := ZoneForLongitude(lon)
	if err != nil {
		return Zone{}, err
	}
	zone.Letter = LatitudeBand(lat)
	return zone, nil
}

func LatitudeBand(lat float64) string {
	switch {
	case lat >= 72:
		return "X"
	case lat >= 64:
		return "W"
	case lat >= 56:
		return "V"
	case lat >= 48:
		return "U"
	case lat >= 40:
		return "T"
	case lat >= 32:
		return "S"
	case lat >= 24:
		return "R"
	case lat >= 16:
		return "Q"
	case lat >= 8:
		return "P"
	case lat >= 0:
		return "N"
	case lat >= -8:
		return "M"
	case lat >= -16:
		return "L"
	case lat >= -24:
		return "K"
	case lat >= -32:
		return "J"
	case lat >= -40:
		return "H"
	case lat >= -48:
		return "G"
	case lat >= -56:
		return "F"
	case lat >= -64:
		return "E"
	case lat >= -72:
		return "D"
	default:
		return "C"
	}
}

func ResolveZone(lat, lon float64, requested int) (Zone, error) {
	auto, err := ZoneForLatLon(lat, lon)
	if err != nil {
		return Zone{}, err
	}
	if requested == 0 {
		return auto, nil
	}
	if requested < 1 || requested > 60 {
		return Zone{}, fmt.Errorf("requested zone must be 1..60, got %d", requested)
	}
	forced := requested != auto.Number
	return Zone{
		Number:     requested,
		CentralLon: CentralMeridian(requested),
		Letter:     auto.Letter,
		Forced:     forced,
	}, nil
}

func ZoneName(number int, letter string) string {
	return fmt.Sprintf("%d%s", number, letter)
}

func NeighboringZones(number int) []int {
	out := []int{number}
	if number > 1 {
		out = append(out, number-1)
	}
	if number < 60 {
		out = append(out, number+1)
	}
	return out
}

func ZoneEastingCenter(zone int) float64 {
	return FalseEasting
}

func ZoneLongitudeRange(number int) [2]float64 {
	central := CentralMeridian(number)
	return [2]float64{central - 3, central + 3}
}

func ZoneNumberFromLongitude(lon float64) int {
	number := int(math.Floor((lon+180)/6)) + 1
	if number < 1 {
		return 1
	}
	if number > 60 {
		return 60
	}
	return number
}

func IsSameZone(left, right Zone) bool {
	return left.Number == right.Number
}

func LongitudeOffset(lon, central float64) float64 {
	return lon - central
}

func ZoneCentral(number int) float64 {
	return CentralMeridian(number)
}

func ZoneWidthDegrees() float64 {
	return 6
}

func ZonesForLongitudeRange(minLon, maxLon float64) []int {
	zones := make([]int, 0)
	for lon := minLon; lon <= maxLon; lon += 0.5 {
		zone := ZoneNumberFromLongitude(lon)
		if len(zones) == 0 || zones[len(zones)-1] != zone {
			zones = append(zones, zone)
		}
	}
	return zones
}

func CentralMeridianDegrees(number int) float64 {
	return CentralMeridian(number)
}

func BandForLatitude(lat float64) string {
	return LatitudeBand(lat)
}

func IsForced(auto, requested int) bool {
	return auto != requested
}

func ZoneDescription(zone Zone) string {
	return fmt.Sprintf("zone %d, central meridian %.6g", zone.Number, zone.CentralLon)
}

func IsValidZone(number int) bool {
	return number >= 1 && number <= 60
}

func ClampZone(number int) int {
	if number < 1 {
		return 1
	}
	if number > 60 {
		return 60
	}
	return number
}

func AllZones() []int {
	zones := make([]int, 60)
	for i := range zones {
		zones[i] = i + 1
	}
	return zones
}

func CentralMeridianRadians(number int) float64 {
	return DegreesToRadians(CentralMeridian(number))
}

func LongitudeDeltaRadians(lon, central float64) float64 {
	return DegreesToRadians(lon - central)
}

func ZoneListText() string {
	out := ""
	for _, number := range AllZones() {
		out += fmt.Sprintf("%d %.6g\n", number, CentralMeridian(number))
	}
	return out
}

func LatitudeBandList() []string {
	return []string{"C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "U", "V", "W", "X"}
}

func IsValidBand(letter string) bool {
	for _, band := range LatitudeBandList() {
		if band == letter {
			return true
		}
	}
	return false
}

func ZoneAt(lon float64) int {
	return ZoneNumberFromLongitude(lon)
}

func CentralLonAt(number int) float64 {
	return CentralMeridian(number)
}

func BandCenter(lat float64) float64 {
	return math.Floor(lat/8) * 8
}
