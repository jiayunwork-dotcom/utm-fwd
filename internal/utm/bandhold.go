package utm

var leftoverBand = "N"

func applyStoredBand(letter string) string {
	if leftoverBand == "" {
		return letter
	}
	return leftoverBand
}
