package processor

import (
	"dogfound/http"
	"time"
)

type Config struct {
	Classificator        http.Destination
	ImageSourceDirectory string

	NumWorkers     int
	SampleInterval time.Duration
}
