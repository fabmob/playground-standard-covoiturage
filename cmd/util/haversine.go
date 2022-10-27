package util

import "github.com/umahmood/haversine"

// Coord stores a position with latitude and longitude data
type Coord haversine.Coord

// Distance returns the distance in kilometers between two positions. Uses the
// haversine formula.
func Distance(coord1, coord2 Coord) float64 {
	_, d := haversine.Distance(haversine.Coord(coord1), haversine.Coord(coord2))
	return d
}

// CoordIgnore can be used as default position when the position does not
// matter
var CoordIgnore = Coord{0, 0}
