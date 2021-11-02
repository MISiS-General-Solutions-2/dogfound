package main

import (
	"dogfound/api"
	"dogfound/database"
	"dogfound/http"
	"dogfound/processor"
	"time"
)

func main() {

	// will get from command line
	// classificatorAddress := "neural_network:6002"
	classificatorAddress := "localhost:6002"
	imageSourceDirectory := database.DataPath + "new_images/"

	close := database.Connect()
	defer close()

	go processor.StartProcessor(&processor.Config{
		Classificator:        http.Destination{Address: classificatorAddress, Retries: 3, RetryInterval: 1 * time.Second},
		ImageSourceDirectory: imageSourceDirectory,
		NumWorkers:           1,
		SampleInterval:       1 * time.Second,
	})

	api.Serve()
}
