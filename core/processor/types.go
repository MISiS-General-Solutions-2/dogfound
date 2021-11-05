package processor

import (
	"dogfound/http"
	"time"
)

type Config struct {
	Classificator           http.Destination
	CameraInputDirectory    string
	VolunteerInputDirectory string

	NumWorkers     int
	SampleInterval time.Duration
}
type volunteerAddedImage struct {
	filename  string
	timestamp int
	lonlat    [2]float64
}
