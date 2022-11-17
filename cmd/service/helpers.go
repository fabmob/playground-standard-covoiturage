package service

// Max returns the larger of x or y.
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
