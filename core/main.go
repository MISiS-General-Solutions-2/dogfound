package main

import (
	"fmt"
	"pet-track/api"
	"pet-track/database"
)

func dbInitMock() {
	files, err := database.GetFilesInDirectory("../data/img")
	if err != nil {
		panic(err)
	}
	if err = database.AddImages(files); err != nil {
		panic(err)
	}
	records := []database.Record{
		{
			Img:   "furry-nyanners.png",
			Color: 1,
			Tail:  1,
		},
		{
			Img:   "defender.png",
			Color: 2,
			Tail:  0,
		},
	}
	if err := database.SetFeatures(records); err != nil {
		panic(err)
	}
	features := map[string]interface{}{
		database.FeatureColor: "0",
	}
	imgs, err := database.GetImagesByFeatures(features)
	if err != nil {
		panic(err)
	}
	fmt.Println(imgs)
}

func main() {
	closer := database.Connect()
	defer closer()

	api.Serve()
}
