package api

import "dogfound/database"

type SimilarResponse struct {
	IsAnimal int                       `json:"is_animal_there"`
	Results  []database.SearchResponse `json:"results"`
}
type PredictRouteRequest struct {
	Lonlat    [2]float64 `json:"lonlat"`
	Timestamp int64      `json:"timestamp"`
}
