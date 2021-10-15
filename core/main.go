package main

import (
	"fmt"
	"pet-track/database"
)

func dbInitMock() {
	files := database.GetImages()
	if err := database.AddImages(files); err != nil {
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
	// closer := database.Connect()
	// defer closer()

	// api.Serve()

	// imgs := []string{
	// 	`C:\Users\antonvlasov\Desktop\img\660_cut.pgm`,
	// 	`C:\Users\antonvlasov\Desktop\img\inverted.pgm`,
	// }
	// dest := `C:\Users\antonvlasov\Desktop\img\merged.pgm`

	database.PopulateWithImages(`D:\Papka\work\MGS2\data\train`)
	// if parsed, err := cv.GetImagesText(imgs); err != nil {
	// 	panic(err)
	// } else {
	// 	for i := range parsed {
	// 		if parsed[i] == "" {
	// 			fmt.Printf("did not parse image %v\n", imgs[i])
	// 		}
	// 	}
	// 	//fmt.Println(parsed)
	// }

}
