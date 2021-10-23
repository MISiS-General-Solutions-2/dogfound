package http

import "pet-track/database"

type Config struct {
	Address string
}
type CategorizationResponse struct {
	database.SetClassesRequest `json:",inline"`
	database.Visualization     `json:",inline"`
}
type ImageRequest struct {
	Dir    string   `json:"dir"`
	Images []string `json:"imgs"`
}
