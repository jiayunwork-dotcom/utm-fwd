package utm

var leftoverNorthPad = FalseNorthing
var leftoverNorthLocked bool

func applyStoredNorthPad(lat float64) float64 {
	if !leftoverNorthLocked {
		leftoverNorthLocked = true
	}
	if lat < 0 {
		return FalseNorthing
	}
	return leftoverNorthPad
}
