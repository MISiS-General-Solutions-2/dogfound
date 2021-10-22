package main

import (
	"pet-track/api"
	"pet-track/database"
	"pet-track/processor"
)

func main() {
	closer := database.Connect()
	defer closer()

	if err := processor.ProcessAllImages(nil); err != nil {
		panic(err)
	}

	api.Serve()
}
