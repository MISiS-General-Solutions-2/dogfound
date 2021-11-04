package main

import (
	"dogfound/api"
	"dogfound/database"
	"dogfound/http"
	"dogfound/processor"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	imageSourceDirectory = database.DataPath + "new_images/"
)

var (
	classificatorAddress string
	numWorkers           int
	sampleInterval       int
)

func simpleParseArgs() {
	if len(os.Args) != 4 {
		panic("must provide arguments: dogfound classificator_address num_workers sample_interval")
	}
	var err error
	classificatorAddress = os.Args[1]
	numWorkers, err = strconv.Atoi(os.Args[2])
	if err != nil {
		panic("num_workers must be integer")
	}
	sampleInterval, err = strconv.Atoi(os.Args[2])
	if err != nil {
		panic("sample_interval must be integer")
	}
}

func main() {
	// unparsed := []string{"B102.jpg", "B105.jpg", "B106.jpg", "B107.jpg", "B11.jpg", "B111.jpg", "B112.jpg", "B116.jpg", "B117.jpg", "B118.jpg", "B119.jpg", "B14.jpg", "B120.jpg", "B19.jpg", "B16.jpg", "B18.jpg", "B23.jpg", "B24.jpg", "B25.jpg", "B26.jpg", "B39.jpg", "B4.jpg", "B44.jpg", "B48.jpg", "B57.jpg", "B56.jpg", "B59.jpg", "B74.jpg", "B78.jpg", "B8.jpg", "B81.jpg", "B82.jpg", "А102.jpg", "А13.jpg", "А18.jpg", "А19.jpg", "А27.jpg", "А31.jpg", "А32.jpg", "А38.jpg", "А39.jpg", "А49.jpg", "А51.jpg", "А55.jpg", "А59.jpg", "А6.jpg", "А60.jpg", "А66.jpg", "А67.jpg", "А7.jpg", "А78.jpg", "А79.jpg", "А8.jpg", "А82.jpg", "А9.jpg", "А90.jpg", "А91.jpg", "А93.jpg", "А98.jpg", "С100.jpg", "С102.jpg", "С103.jpg", "С21.jpg", "С25.jpg", "С34.jpg", "С35.jpg", "С38.jpg", "С44.jpg", "С47.jpg", "С48.jpg", "С54.jpg", "С62.jpg", "С70.jpg", "С78.jpg", "С94.jpg"}
	// unmatched := []string{"B37.jpg", "B75.jpg", "А10.jpg", "А101.jpg", "А12.jpg", "А16.jpg", "А15.jpg", "А11.jpg", "А17.jpg", "А14.jpg", "А21.jpg", "А20.jpg", "А23.jpg", "А24.jpg", "А22.jpg", "А29.jpg", "А26.jpg", "А28.jpg", "А25.jpg", "А30.jpg", "А40.jpg", "А41.jpg", "А42.jpg", "А43.jpg", "А44.jpg", "А45.jpg", "А46.jpg", "А47.jpg", "А50.jpg", "А52.jpg", "А81.jpg", "А83.jpg", "С15.jpg", "С16.jpg", "С27.jpg", "С3.jpg", "С39.jpg", "С41.jpg", "С51.jpg", "С61.jpg", "С65.jpg", "С74.jpg", "С79.jpg", "С86.jpg", "С89.jpg", "С85.jpg", "С9.jpg", "С98.jpg"}
	// srcDir := "/home/vlasov/folder/mgs2/dogfound/core/my_data/new_images/"
	// unparsedDdt := "/home/vlasov/folder/mgs2/dogfound/core/my_data/unparsed/"
	// unmatchedDst := "/home/vlasov/folder/mgs2/dogfound/core/my_data/unmatched/"
	// if err := cv.MoveSelectedToFile(srcDir, unparsedDdt, unparsed); err != nil {
	// 	panic(err)
	// }
	// if err := cv.MoveSelectedToFile(srcDir, unmatchedDst, unmatched); err != nil {
	// 	panic(err)
	// }
	simpleParseArgs()
	fmt.Println(classificatorAddress)

	close := database.Connect()
	defer close()

	go processor.StartProcessor(&processor.Config{
		Classificator:        http.Destination{Address: classificatorAddress, Retries: 3, RetryInterval: 1 * time.Second},
		ImageSourceDirectory: imageSourceDirectory,
		NumWorkers:           numWorkers,
		SampleInterval:       time.Duration(sampleInterval) * time.Second,
	})

	api.Serve()
}
