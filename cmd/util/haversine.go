package util

import "github.com/umahmood/haversine"

type Coord haversine.Coord

func Distance(coord1, coord2 Coord) float64 {
	_, d := haversine.Distance(haversine.Coord(coord1), haversine.Coord(coord2))
	return d
}

var CoordIgnore = Coord{0, 0}
