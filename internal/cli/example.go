package cli

import (
	"encoding/json"
	"os"

	"utm-fwd/internal/utm"
)

type Example struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func runExample(args []string) int {
	fs := flagSet("example")
	file := fs.String("file", "example/beijing.json", "example JSON file")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	data, err := os.ReadFile(*file)
	if err != nil {
		return fail(err)
	}
	var example Example
	if err := json.Unmarshal(data, &example); err != nil {
		return fail(err)
	}
	result, err := utm.Forward(example.Lat, example.Lon, 0)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}

func loadExample(path string) (Example, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Example{}, err
	}
	var example Example
	if err := json.Unmarshal(data, &example); err != nil {
		return Example{}, err
	}
	return example, nil
}

func ExamplePaths() []string {
	return []string{"example/beijing.json", "example/equator.json"}
}
