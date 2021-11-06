package geo

import (
	"dogfound/database"
	"math"
	"time"
)

const (
	cosineMoscow            = 0.562576257
	dogSearchRadiusMeters   = 5000
	dogSearchTimeframeHours = 72
)

func GetPossibleRoute(timestamp time.Time, lonlat [2]float64) ([]database.SearchResponse, error) {
	lonMeters, latMeters := lonToMetersApprox(lonlat[0]), latToMetersApprox(lonlat[1])
	var lonLatBox [2][2]float64
	lonLatBox[0][0] = metersToLonApprox(lonMeters - dogSearchRadiusMeters)
	lonLatBox[0][1] = metersToLatApprox(latMeters - dogSearchRadiusMeters)
	lonLatBox[1][0] = metersToLonApprox(lonMeters + dogSearchRadiusMeters)
	lonLatBox[1][1] = metersToLatApprox(latMeters + dogSearchRadiusMeters)

	t0 := timestamp.Add(-time.Duration(dogSearchTimeframeHours) * time.Hour)
	t1 := timestamp.Add(time.Duration(dogSearchTimeframeHours) * time.Hour)

	return database.GetImagesWithinFrame(lonLatBox, t0.Unix(), t1.Unix())
}

func metersToLatApprox(dist float64) float64 {
	return dist / 111111
}
func metersToLonApprox(dist float64) float64 {
	return dist / (111111 * cosineMoscow)
}
func latToMetersApprox(lat float64) float64 {
	return lat * 111111
}
func lonToMetersApprox(lat float64) float64 {
	return lat * 111111 * cosineMoscow
}
func euclidianDistanceApprox(aLat, bLat, aLng, bLng float64) float64 {
	dy := 111111 * (aLat - bLat)
	dx := 111111 * (aLng - bLng) * cosineMoscow
	return math.Sqrt(dx*dx + dy*dy)
}
