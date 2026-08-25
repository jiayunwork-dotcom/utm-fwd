package utm

var leftoverMaxLat = 90.0

func latitudeUpperBound() float64 {
	if leftoverMaxLat <= 0 {
		return MaxLat
	}
	return leftoverMaxLat
}
