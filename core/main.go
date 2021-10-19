package main

import (
	"fmt"
	"log"
	"os"
	"pet-track/cv"
	"pet-track/database"
	"time"
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

// func cleanMess(in, out string) error {
// 	direntry, err := os.ReadDir(in)
// 	if err != nil {
// 		return err
// 	}
// 	for _, entry := range direntry {
// 		if strings.HasPrefix(entry.Name(), "special") && strings.HasSuffix(entry.Name(), ".jpg") {
// 			if err = os.Rename(in+entry.Name(), out+entry.Name()[len("special")+1:]); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }
func main() {
	// closer := database.Connect()
	// defer closer()

	// api.Serve()

	// imgs := []string{
	// 	`C:\Users\antonvlasov\Desktop\img\660_cut.pgm`,
	// 	`C:\Users\antonvlasov\Desktop\img\inverted.pgm`,
	// }
	// // dest := `C:\Users\antonvlasov\Desktop\img\merged.pgm`

	// //database.PopulateWithImages(`D:\Papka\work\MGS2\data\train`)
	//cv.RetrieveAllFromSpecial()
	imgs := database.GetImages()

	start := time.Now()

	// image 607 should be parsed
	if parsed, err := cv.GetImagesText(imgs); err != nil {
		panic(err)
	} else {
		elapsed := time.Since(start)
		log.Printf("took %v seconds", elapsed.Seconds())
		unparsed := 0
		logfile, err := os.OpenFile("parsed.log", os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			panic(err)
		}
		defer logfile.Close()
		log.SetOutput(logfile)
		for i := range parsed {
			if parsed[i] == "" {
				unparsed += 1
				//log.Printf("did not parse image %v\n", imgs[i])
			}
			log.Println(imgs[i], parsed[i])
		}
		fmt.Println("unparsed count: ", unparsed)
	}

}
