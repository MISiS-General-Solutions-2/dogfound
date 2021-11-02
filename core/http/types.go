package http

import (
	"dogfound/database"
	"time"
)

type CategorizationResponse struct {
	database.ClassInfo `json:",inline"`
	Vis                database.Visualization `json:"vis"`
}
type ImageRequest struct {
	Image string `json:"image"`
}
type Destination struct {
	Address       string
	Retries       int
	RetryInterval time.Duration
}
