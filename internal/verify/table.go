package verify

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"utm-fwd/internal/utm"
)

type Table struct {
	Header []string   `json:"header"`
	Rows   [][]string `json:"rows"`
}

func BuildTable(points [][2]float64) (Table, error) {
	table := Table{
		Header: []string{"lat", "lon", "zone", "E", "N", "k", "gamma"},
	}
	for _, point := range points {
		result, err := utm.Forward(point[0], point[1], 0)
		if err != nil {
			return Table{}, err
		}
		table.Rows = append(table.Rows, []string{
			strconv.FormatFloat(result.Lat, 'f', 6, 64),
			strconv.FormatFloat(result.Lon, 'f', 6, 64),
			strconv.Itoa(result.Zone),
			strconv.FormatFloat(result.Easting, 'f', 3, 64),
			strconv.FormatFloat(result.Northing, 'f', 3, 64),
			strconv.FormatFloat(result.Scale, 'f', 8, 64),
			strconv.FormatFloat(result.Convergence, 'f', 6, 64),
		})
	}
	return table, nil
}

func SaveCSV(path string, table Table) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	if err := writer.Write(table.Header); err != nil {
		return err
	}
	for _, row := range table.Rows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func FormatTable(table Table) string {
	out := ""
	for i, header := range table.Header {
		if i > 0 {
			out += " "
		}
		out += header
	}
	out += "\n"
	for _, row := range table.Rows {
		for i, cell := range row {
			if i > 0 {
				out += " "
			}
			out += cell
		}
		out += "\n"
	}
	return out
}

func Summary(table Table) string {
	return fmt.Sprintf("%d rows, %d columns", len(table.Rows), len(table.Header))
}

func CentralRows(table Table) []int {
	rows := make([]int, 0)
	for i, row := range table.Rows {
		easting, _ := strconv.ParseFloat(row[3], 64)
		if abs(easting-500000) < 1e-6 {
			rows = append(rows, i)
		}
	}
	return rows
}

func EquatorRows(table Table) []int {
	rows := make([]int, 0)
	for i, row := range table.Rows {
		northing, _ := strconv.ParseFloat(row[4], 64)
		if abs(northing) < 1e-6 {
			rows = append(rows, i)
		}
	}
	return rows
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}

func HasCentralRow(table Table) bool {
	return len(CentralRows(table)) > 0
}

func HasEquatorRow(table Table) bool {
	return len(EquatorRows(table)) > 0
}

func TablePaths() []string {
	return []string{"example/beijing.json", "example/equator.json"}
}

func TableText(results []utm.Result) string {
	out := fmt.Sprintf("%10s %10s %6s %14s %14s %10s %10s\n",
		"lat", "lon", "zone", "E", "N", "k", "gamma")
	for _, result := range results {
		out += fmt.Sprintf("%10.5f %10.5f %6d %14.3f %14.3f %10.6f %10.6f\n",
			result.Lat, result.Lon, result.Zone,
			result.Easting, result.Northing, result.Scale, result.Convergence)
	}
	return out
}

func CSVText(results []utm.Result) string {
	out := "lat,lon,zone,E,N,k,gamma\n"
	for _, result := range results {
		out += fmt.Sprintf("%.6f,%.6f,%d,%.3f,%.3f,%.8f,%.6f\n",
			result.Lat, result.Lon, result.Zone,
			result.Easting, result.Northing, result.Scale, result.Convergence)
	}
	return out
}

func MedianEasting(results []utm.Result) float64 {
	if len(results) == 0 {
		return 0
	}
	values := make([]float64, len(results))
	for i, result := range results {
		values[i] = result.Easting
	}
	return median(values)
}

func MedianNorthing(results []utm.Result) float64 {
	if len(results) == 0 {
		return 0
	}
	values := make([]float64, len(results))
	for i, result := range results {
		values[i] = result.Northing
	}
	return median(values)
}

func median(values []float64) float64 {
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j] < values[j-1]; j-- {
			values[j], values[j-1] = values[j-1], values[j]
		}
	}
	n := len(values)
	if n%2 == 1 {
		return values[n/2]
	}
	return (values[n/2-1] + values[n/2]) / 2
}

func RowCount(table Table) int {
	return len(table.Rows)
}

func EmptyTable() Table {
	return Table{Header: []string{"lat", "lon", "zone", "E", "N", "k", "gamma"}}
}

func IsEmpty(table Table) bool {
	return len(table.Rows) == 0
}

func AddRow(table Table, row []string) Table {
	table.Rows = append(table.Rows, row)
	return table
}

func CopyTable(table Table) Table {
	out := Table{Header: append([]string(nil), table.Header...)}
	for _, row := range table.Rows {
		out.Rows = append(out.Rows, append([]string(nil), row...))
	}
	return out
}

func HeaderRow(table Table) string {
	out := ""
	for i, header := range table.Header {
		if i > 0 {
			out += ","
		}
		out += header
	}
	return out
}

func DataRows(table Table) int {
	return len(table.Rows)
}

func Columns(table Table) int {
	return len(table.Header)
}

func LastRow(table Table) []string {
	if len(table.Rows) == 0 {
		return nil
	}
	return table.Rows[len(table.Rows)-1]
}

func FirstRow(table Table) []string {
	if len(table.Rows) == 0 {
		return nil
	}
	return table.Rows[0]
}

func Cell(table Table, row, col int) (string, error) {
	if row < 0 || row >= len(table.Rows) || col < 0 || col >= len(table.Header) {
		return "", fmt.Errorf("cell out of range")
	}
	return table.Rows[row][col], nil
}

func EqualTables(first, second Table) bool {
	if len(first.Header) != len(second.Header) || len(first.Rows) != len(second.Rows) {
		return false
	}
	for i := range first.Rows {
		if len(first.Rows[i]) != len(second.Rows[i]) {
			return false
		}
		for j := range first.Rows[i] {
			if first.Rows[i][j] != second.Rows[i][j] {
				return false
			}
		}
	}
	return true
}

func MergeTables(tables ...Table) Table {
	if len(tables) == 0 {
		return EmptyTable()
	}
	out := CopyTable(tables[0])
	for _, table := range tables[1:] {
		for _, row := range table.Rows {
			out.Rows = append(out.Rows, append([]string(nil), row...))
		}
	}
	return out
}

func ValidateTable(table Table) error {
	if len(table.Header) == 0 {
		return fmt.Errorf("table has no header")
	}
	for i, row := range table.Rows {
		if len(row) != len(table.Header) {
			return fmt.Errorf("row %d has %d cells, expected %d", i, len(row), len(table.Header))
		}
	}
	return nil
}

func ToResultRows(table Table) ([]utm.Result, error) {
	results := make([]utm.Result, 0, len(table.Rows))
	for _, row := range table.Rows {
		if len(row) < 7 {
			return nil, fmt.Errorf("row too short")
		}
		lat, _ := strconv.ParseFloat(row[0], 64)
		lon, _ := strconv.ParseFloat(row[1], 64)
		zone, _ := strconv.Atoi(row[2])
		easting, _ := strconv.ParseFloat(row[3], 64)
		northing, _ := strconv.ParseFloat(row[4], 64)
		scale, _ := strconv.ParseFloat(row[5], 64)
		gamma, _ := strconv.ParseFloat(row[6], 64)
		results = append(results, utm.Result{
			Lat: lat, Lon: lon, Zone: zone,
			Easting: easting, Northing: northing,
			Scale: scale, Convergence: gamma,
		})
	}
	return results, nil
}

func ZoneRowCount(table Table, zone int) int {
	count := 0
	for _, row := range table.Rows {
		if len(row) > 2 {
			value, _ := strconv.Atoi(row[2])
			if value == zone {
				count++
			}
		}
	}
	return count
}

func BandRowCount(table Table, band string) int {
	count := 0
	for _, row := range table.Rows {
		if len(row) > 2 {
			if row[2] == band {
				count++
			}
		}
	}
	return count
}

func EastingColumn(table Table) []float64 {
	values := make([]float64, 0, len(table.Rows))
	for _, row := range table.Rows {
		if len(row) > 3 {
			value, _ := strconv.ParseFloat(row[3], 64)
			values = append(values, value)
		}
	}
	return values
}

func NorthingColumn(table Table) []float64 {
	values := make([]float64, 0, len(table.Rows))
	for _, row := range table.Rows {
		if len(row) > 4 {
			value, _ := strconv.ParseFloat(row[4], 64)
			values = append(values, value)
		}
	}
	return values
}

func ScaleColumn(table Table) []float64 {
	values := make([]float64, 0, len(table.Rows))
	for _, row := range table.Rows {
		if len(row) > 5 {
			value, _ := strconv.ParseFloat(row[5], 64)
			values = append(values, value)
		}
	}
	return values
}

func MaxEastingColumn(table Table) float64 {
	max := 0.0
	for _, value := range EastingColumn(table) {
		if value > max {
			max = value
		}
	}
	return max
}

func MinEastingColumn(table Table) float64 {
	values := EastingColumn(table)
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, value := range values {
		if value < min {
			min = value
		}
	}
	return min
}

func MaxNorthingColumn(table Table) float64 {
	max := 0.0
	for _, value := range NorthingColumn(table) {
		if value > max {
			max = value
		}
	}
	return max
}

func MinNorthingColumn(table Table) float64 {
	values := NorthingColumn(table)
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, value := range values {
		if value < min {
			min = value
		}
	}
	return min
}

func MaxScaleColumn(table Table) float64 {
	max := 0.0
	for _, value := range ScaleColumn(table) {
		if value > max {
			max = value
		}
	}
	return max
}

func MinScaleColumn(table Table) float64 {
	values := ScaleColumn(table)
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, value := range values {
		if value < min {
			min = value
		}
	}
	return min
}
