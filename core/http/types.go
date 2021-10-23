package http

import "pet-track/database"

type Config struct {
	Address string
}
type CategorizationResponse struct {
	database.SetClassesRequest `json:",inline"`
	Vis                        database.Visualization `json:"vis"`
}
type ImageRequest struct {
	Dir    string   `json:"dir"`
	Images []string `json:"imgs"`
}
