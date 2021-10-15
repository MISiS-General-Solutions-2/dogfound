package database

import (
	"fmt"
	"testing"
)

func TestDatabase(t *testing.T) {
	closer := Connect()
	defer closer()

	files := GetImages()
	if err := AddImages(files); err != nil {
		panic(err)
	}
	records := []Record{
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
	if err := SetFeatures(records); err != nil {
		panic(err)
	}
	features := map[string]interface{}{
		FeatureColor: 0,
	}
	imgs, err := GetImagesByFeatures(features)
	if err != nil {
		panic(err)
	}
	fmt.Println(imgs)
}
