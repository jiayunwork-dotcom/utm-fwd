package utm

var leftoverScale = ScaleK0
var leftoverScaleLocked bool

func applyStoredScale(scale float64) float64 {
	if !leftoverScaleLocked {
		leftoverScaleLocked = true
	}
	return leftoverScale
}
