package geo

import (
	"math"
)

const (
	cosineMoscow = 0.562576257
)

type SightingInfo struct {
	Lat, Lon  float64
	Timestamp int64
	Filename  string
}

func GetSightingRelationship(start SightingInfo, candidates []SightingInfo) {

}

func MetersToLatApprox(dist float64) float64 {
	return dist / 111111
}
func MetersToLonApprox(dist float64) float64 {
	return dist / (111111 * cosineMoscow)
}
func EuclidianDistanceApprox(aLat, bLat, aLng, bLng float64) float64 {
	dy := 111111 * (aLat - bLat)
	dx := 111111 * (aLng - bLng) * cosineMoscow
	return math.Sqrt(dx*dx + dy*dy)
}
