package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"pet-track/cv"
	"pet-track/database"
	"strconv"
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
	imgs := database.GetImages()

	start := time.Now()

	if parsed, err := cv.GetImagesText(imgs); err != nil {
		panic(err)
	} else {
		elapsed := time.Since(start)
		log.Printf("took %v seconds", elapsed.Seconds())
		unparsed := 0
		for i := range parsed {
			if parsed[i] == "" {
				unparsed += 1
				//log.Printf("did not parse image %v\n", imgs[i])
			}
		}
		fmt.Println(unparsed)
		//fmt.Println(parsed)
	}

	// img, err := os.ReadFile(`C:\Users\antonvlasov\Desktop\img\inverted.pgm`)
	// if err != nil {
	// 	panic(err)
	// }

	// c := cv.CallCLibMain(img, []string{"gocr", "-C", "0-9a-zA-Z"})
	// fmt.Println(c)
}
func testCRead(imgfile, out string) error {
	img, err := os.ReadFile(imgfile)
	if err != nil {
		panic(err)
	}
	log, err := os.Open(out)
	if err != nil {
		panic(err)
	}
	s := bufio.NewScanner(log)
	i := 0
	for s.Scan() {
		line := s.Text()
		logByte, err := strconv.Atoi(line)
		if err != nil {
			return err
		}
		if i >= len(img) {
			return fmt.Errorf("too many values at log file")
		}
		lineByte := img[i]
		if logByte != int(lineByte) {
			return fmt.Errorf("different byte at pos %v: expected %v, got %v", i, img[i], line)
		}
		i += 1
	}
	if i != len(img) {
		fmt.Println("different file lengths")
	}
	return nil
}
