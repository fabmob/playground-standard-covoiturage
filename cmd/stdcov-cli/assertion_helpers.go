package main

import "github.com/umahmood/haversine"

type coords struct {
	lat float64
	lon float64
}

func distanceKm(coords1, coords2 coords) float64 {
	c1 := haversine.Coord{Lat: coords1.lat, Lon: coords1.lon}
	c2 := haversine.Coord{Lat: coords2.lat, Lon: coords2.lon}
	_, dist := haversine.Distance(c1, c2)
	return dist
}
