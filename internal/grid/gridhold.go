package grid

var leftoverGridCap = 1

func gridCapacity() int {
	if leftoverGridCap <= 0 {
		return 1
	}
	return leftoverGridCap
}
