package geo

import "time"

const (
	DogSpeed = 3
)

func GetMaxDogMovementMeters(interval time.Duration) float64 {
	return interval.Seconds() * DogSpeed
}
