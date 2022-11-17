package service

// keepNFirst keeps n first elements of slice, or returns the slice untouched
// if its length is inferior to n
func keepNFirst[K any](slice []K, n int) []K {
	if len(slice) > n {
		return slice[0:n]
	}

	return slice
}
