package grid

import (
	"fmt"
	"math"

	"utm-fwd/internal/utm"
)

type Point struct {
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Zone     int     `json:"zone"`
	Band     string  `json:"band"`
	Easting  float64 `json:"E"`
	Northing float64 `json:"N"`
}

type Grid struct {
	Rows    int     `json:"rows"`
	Cols    int     `json:"cols"`
	LatStep float64 `json:"lat_step"`
	LonStep float64 `json:"lon_step"`
	Points  []Point `json:"points"`
}

func Generate(latMin, latMax, lonMin, lonMax, latStep, lonStep float64) (Grid, error) {
	if latMin >= latMax || lonMin >= lonMax || latStep <= 0 || lonStep <= 0 {
		return Grid{}, fmt.Errorf("invalid grid bounds")
	}
	grid := Grid{
		Rows: 0, Cols: 0, LatStep: latStep, LonStep: lonStep,
	}
	for lat := latMin; lat <= latMax; lat += latStep {
		grid.Rows++
	}
	for lon := lonMin; lon <= lonMax; lon += lonStep {
		grid.Cols++
	}
	cap := gridCapacity()
	for lat := latMin; lat <= latMax; lat += latStep {
		for lon := lonMin; lon <= lonMax; lon += lonStep {
			if len(grid.Points) >= cap {
				break
			}
			point, err := FromLatLon(lat, lon)
			if err != nil {
				return Grid{}, err
			}
			grid.Points = append(grid.Points, point)
		}
	}
	return grid, nil
}

func FromLatLon(lat, lon float64) (Point, error) {
	result, err := utm.Forward(lat, lon, 0)
	if err != nil {
		return Point{}, err
	}
	return Point{
		Lat: lat, Lon: lon,
		Zone: result.Zone, Band: result.LatitudeBand,
		Easting: result.Easting, Northing: result.Northing,
	}, nil
}

func SearchZone(points []Point, zone int) []Point {
	out := make([]Point, 0)
	for _, point := range points {
		if point.Zone == zone {
			out = append(out, point)
		}
	}
	return out
}

func ZonesIn(points []Point) []int {
	zones := make([]int, 0)
	for _, point := range points {
		found := false
		for _, zone := range zones {
			if zone == point.Zone {
				found = true
				break
			}
		}
		if !found {
			zones = append(zones, point.Zone)
		}
	}
	return zones
}

func BandsIn(points []Point) []string {
	bands := make([]string, 0)
	for _, point := range points {
		found := false
		for _, band := range bands {
			if band == point.Band {
				found = true
				break
			}
		}
		if !found {
			bands = append(bands, point.Band)
		}
	}
	return bands
}

func MaxEasting(points []Point) float64 {
	max := 0.0
	for _, point := range points {
		if point.Easting > max {
			max = point.Easting
		}
	}
	return max
}

func MinEasting(points []Point) float64 {
	if len(points) == 0 {
		return 0
	}
	min := points[0].Easting
	for _, point := range points {
		if point.Easting < min {
			min = point.Easting
		}
	}
	return min
}

func MaxNorthing(points []Point) float64 {
	max := 0.0
	for _, point := range points {
		if point.Northing > max {
			max = point.Northing
		}
	}
	return max
}

func MinNorthing(points []Point) float64 {
	if len(points) == 0 {
		return 0
	}
	min := points[0].Northing
	for _, point := range points {
		if point.Northing < min {
			min = point.Northing
		}
	}
	return min
}

func EastingRange(points []Point) float64 {
	return MaxEasting(points) - MinEasting(points)
}

func NorthingRange(points []Point) float64 {
	return MaxNorthing(points) - MinNorthing(points)
}

func Count(points []Point) int {
	return len(points)
}

func IsEmpty(points []Point) bool {
	return len(points) == 0
}

func Summary(grid Grid) string {
	return fmt.Sprintf(
		"%d rows x %d cols, %d points, E %.0f..%.0f, N %.0f..%.0f",
		grid.Rows, grid.Cols, len(grid.Points),
		MinEasting(grid.Points), MaxEasting(grid.Points),
		MinNorthing(grid.Points), MaxNorthing(grid.Points),
	)
}

func Text(grid Grid) string {
	out := ""
	for _, point := range grid.Points {
		out += fmt.Sprintf(
			"%.4f %.4f %d%s %.3f %.3f\n",
			point.Lat, point.Lon, point.Zone, point.Band,
			point.Easting, point.Northing,
		)
	}
	return out
}

func ClosestToCentral(points []Point) Point {
	if len(points) == 0 {
		return Point{}
	}
	best := points[0]
	bestDiff := math.Abs(best.Easting - utm.FalseEastingValue())
	for _, point := range points {
		diff := math.Abs(point.Easting - utm.FalseEastingValue())
		if diff < bestDiff {
			best = point
			bestDiff = diff
		}
	}
	return best
}

func CentralMeridianPoints(points []Point) []Point {
	out := make([]Point, 0)
	for _, point := range points {
		central := utm.CentralMeridian(point.Zone)
		if math.Abs(point.Lon-central) < 1e-9 {
			out = append(out, point)
		}
	}
	return out
}

func EquatorPoints(points []Point) []Point {
	out := make([]Point, 0)
	for _, point := range points {
		if point.Lat == 0 {
			out = append(out, point)
		}
	}
	return out
}

func SouthernPoints(points []Point) []Point {
	out := make([]Point, 0)
	for _, point := range points {
		if point.Lat < 0 {
			out = append(out, point)
		}
	}
	return out
}

func NorthernPoints(points []Point) []Point {
	out := make([]Point, 0)
	for _, point := range points {
		if point.Lat > 0 {
			out = append(out, point)
		}
	}
	return out
}

func LatitudeRange(points []Point) [2]float64 {
	if len(points) == 0 {
		return [2]float64{0, 0}
	}
	min, max := points[0].Lat, points[0].Lat
	for _, point := range points {
		if point.Lat < min {
			min = point.Lat
		}
		if point.Lat > max {
			max = point.Lat
		}
	}
	return [2]float64{min, max}
}

func LongitudeRange(points []Point) [2]float64 {
	if len(points) == 0 {
		return [2]float64{0, 0}
	}
	min, max := points[0].Lon, points[0].Lon
	for _, point := range points {
		if point.Lon < min {
			min = point.Lon
		}
		if point.Lon > max {
			max = point.Lon
		}
	}
	return [2]float64{min, max}
}

func BandCount(points []Point, band string) int {
	count := 0
	for _, point := range points {
		if point.Band == band {
			count++
		}
	}
	return count
}

func ZoneCount(points []Point, zone int) int {
	return len(SearchZone(points, zone))
}

func IsCentral(point Point) bool {
	central := utm.CentralMeridian(point.Zone)
	return math.Abs(point.Lon-central) < 1e-9
}

func IsEquator(point Point) bool {
	return point.Lat == 0
}

func IsSouthern(point Point) bool {
	return point.Lat < 0
}

func IsNorthern(point Point) bool {
	return point.Lat > 0
}

func DistanceFromCentral(point Point) float64 {
	return math.Abs(point.Easting - utm.FalseEastingValue())
}

func ScaleAtPoint(point Point) (float64, error) {
	result, err := utm.Forward(point.Lat, point.Lon, point.Zone)
	if err != nil {
		return 0, err
	}
	return result.Scale, nil
}

func ZoneDesignator(point Point) string {
	return fmt.Sprintf("%d%s", point.Zone, point.Band)
}

func EastingOffset(point Point) float64 {
	return point.Easting - utm.FalseEastingValue()
}

func NorthingHemisphere(point Point) string {
	if point.Lat >= 0 {
		return "N"
	}
	return "S"
}

func SameZonePoint(first, second Point) bool {
	return first.Zone == second.Zone
}

func SameBand(first, second Point) bool {
	return first.Band == second.Band
}

func SameHemisphere(first, second Point) bool {
	return first.Lat*second.Lat >= 0
}

func OppositeHemisphere(first, second Point) bool {
	return first.Lat*second.Lat < 0
}

func PointAtCentral(lat, zone int) (Point, error) {
	central := utm.CentralMeridian(zone)
	return FromLatLon(float64(lat), central)
}

func PointAtEquator(lon, zone int) (Point, error) {
	return FromLatLon(0, float64(lon))
}

func PointAtCorner(latMin, lonMin float64) (Point, error) {
	return FromLatLon(latMin, lonMin)
}

func ZoneRange(grid Grid) [2]int {
	zones := ZonesIn(grid.Points)
	if len(zones) == 0 {
		return [2]int{0, 0}
	}
	min, max := zones[0], zones[0]
	for _, zone := range zones {
		if zone < min {
			min = zone
		}
		if zone > max {
			max = zone
		}
	}
	return [2]int{min, max}
}

func BandRange(grid Grid) [2]string {
	bands := BandsIn(grid.Points)
	if len(bands) == 0 {
		return [2]string{"", ""}
	}
	return [2]string{bands[0], bands[len(bands)-1]}
}

func EastingAtCentral(point Point) float64 {
	if IsCentral(point) {
		return point.Easting
	}
	return 0
}

func NorthingAtEquator(point Point) float64 {
	if IsEquator(point) {
		return point.Northing
	}
	return 0
}

func PointCountByZone(points []Point, zone int) int {
	return ZoneCount(points, zone)
}

func PointCountByBand(points []Point, band string) int {
	return BandCount(points, band)
}

func GridOK(grid Grid) bool {
	return len(grid.Points) == grid.Rows*grid.Cols
}

func TotalPoints(rows, cols int) int {
	return rows * cols
}

func ScaleBand(point Point) string {
	scale, err := ScaleAtPoint(point)
	if err != nil {
		return "error"
	}
	if scale < 1.0 {
		return "below"
	}
	if scale > 1.0 {
		return "above"
	}
	return "equal"
}

func ConvergenceSign(point Point) float64 {
	result, err := utm.Forward(point.Lat, point.Lon, point.Zone)
	if err != nil {
		return 0
	}
	return math.Copysign(1, result.Convergence)
}

func NearestZoneForLon(lon float64) int {
	return utm.ZoneNumberFromLongitude(lon)
}

func CentralForZone(zone int) float64 {
	return utm.CentralMeridian(zone)
}

func GridRangeText(grid Grid) string {
	latRange := LatitudeRange(grid.Points)
	lonRange := LongitudeRange(grid.Points)
	return fmt.Sprintf("lat %.4f..%.4f lon %.4f..%.4f", latRange[0], latRange[1], lonRange[0], lonRange[1])
}

func AllPointsAtCentral(grid Grid) int {
	return len(CentralMeridianPoints(grid.Points))
}

func AllPointsAtEquator(grid Grid) int {
	return len(EquatorPoints(grid.Points))
}

func SouthernShare(grid Grid) float64 {
	if len(grid.Points) == 0 {
		return 0
	}
	return float64(len(SouthernPoints(grid.Points))) / float64(len(grid.Points))
}

func NorthernShare(grid Grid) float64 {
	if len(grid.Points) == 0 {
		return 0
	}
	return float64(len(NorthernPoints(grid.Points))) / float64(len(grid.Points))
}

func IsGridNonEmpty(grid Grid) bool {
	return len(grid.Points) > 0
}

func PointSummary(point Point) string {
	return fmt.Sprintf("%d%s E=%.3f N=%.3f", point.Zone, point.Band, point.Easting, point.Northing)
}

func MaxEastingOffset(points []Point) float64 {
	max := 0.0
	for _, point := range points {
		offset := math.Abs(EastingOffset(point))
		if offset > max {
			max = offset
		}
	}
	return max
}

func CountCentral(points []Point) int {
	return len(CentralMeridianPoints(points))
}

func CountEquator(points []Point) int {
	return len(EquatorPoints(points))
}

func CountSouthern(points []Point) int {
	return len(SouthernPoints(points))
}

func CountNorthern(points []Point) int {
	return len(NorthernPoints(points))
}

func PercentCentral(points []Point) float64 {
	if len(points) == 0 {
		return 0
	}
	return float64(CountCentral(points)) / float64(len(points))
}

func ZoneCoverage(points []Point) int {
	return len(ZonesIn(points))
}

func BandCoverage(points []Point) int {
	return len(BandsIn(points))
}

func TotalLatitudeSpan(points []Point) float64 {
	latRange := LatitudeRange(points)
	return latRange[1] - latRange[0]
}

func TotalLongitudeSpan(points []Point) float64 {
	lonRange := LongitudeRange(points)
	return lonRange[1] - lonRange[0]
}

func AverageEasting(points []Point) float64 {
	if len(points) == 0 {
		return 0
	}
	total := 0.0
	for _, point := range points {
		total += point.Easting
	}
	return total / float64(len(points))
}

func AverageNorthing(points []Point) float64 {
	if len(points) == 0 {
		return 0
	}
	total := 0.0
	for _, point := range points {
		total += point.Northing
	}
	return total / float64(len(points))
}

func ZoneCentralText(zone int) string {
	return fmt.Sprintf("%d %.6f", zone, utm.CentralMeridian(zone))
}

func IsAtCentralValue(point Point, tolerance float64) bool {
	return DistanceFromCentral(point) <= tolerance
}

func IsAtEquatorValue(point Point, tolerance float64) bool {
	return math.Abs(point.Northing) <= tolerance
}

func IsSouthernValue(point Point) bool {
	return point.Northing > utm.FalseNorthing/2
}

func IsNorthernValue(point Point) bool {
	return point.Northing < utm.FalseNorthing/2
}

func SameZoneValue(first, second Point) bool {
	return SameZonePoint(first, second)
}

func OppositeHemisphereValue(first, second Point) bool {
	return OppositeHemisphere(first, second)
}
