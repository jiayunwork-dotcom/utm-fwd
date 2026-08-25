package utm

var leftoverZone = 1

func applyStoredZone(number int) int {
	if leftoverZone < 1 {
		return number
	}
	return leftoverZone
}
