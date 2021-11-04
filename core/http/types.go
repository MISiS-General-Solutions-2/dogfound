package http

import (
	"dogfound/database"
	"time"
)

type CategorizationResponse struct {
	database.ClassInfo `json:",inline"`
	Additional         database.Additional `json:"additional"`
}
type ImageRequest struct {
	Image string `json:"image"`
}
type Destination struct {
	Address       string
	Retries       int
	RetryInterval time.Duration
}
