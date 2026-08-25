package cli

import (
	"fmt"
	"strconv"

	"utm-fwd/internal/utm"
)

func parseFloat(text, label string) (float64, error) {
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number, got %q", label, text)
	}
	return value, nil
}

func validateLatLon(lat, lon float64) error {
	return utm.ValidateLatLon(lat, lon)
}

func validateZone(zone int) error {
	if zone == 0 {
		return nil
	}
	return utm.ValidateZoneNumber(zone)
}

func requireFinite(value float64, label string) error {
	return utm.RequireFinite(value, label)
}

func normalize(text string) string {
	if len(text) > 0 && text[0] == '"' {
		if parsed, err := strconv.Unquote(text); err == nil {
			return parsed
		}
	}
	return text
}
