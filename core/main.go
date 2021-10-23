package main

import (
	"pet-track/api"
	"pet-track/database"
	"pet-track/processor"
)

func main() {
	database.PopulateWithImages("/home/vlasov/folder/mgs2/pet-track/core/data/img.test")
	closer := database.Connect()
	defer closer()

	if err := processor.ProcessAllImages(nil); err != nil {
		panic(err)
	}

	api.Serve()
}
