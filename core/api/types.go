package api

import "dogfound/database"

type SimilarResponse struct {
	IsAnimal int                       `json:"is_animal_there"`
	Results  []database.SearchResponse `json:"results"`
}
