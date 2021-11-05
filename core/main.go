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
	// unparsed := []string{"B1.jpg", "B102.jpg", "B103.jpg", "B104.jpg", "B105.jpg", "B106.jpg", "B107.jpg", "B108.jpg", "B109.jpg", "B11.jpg", "B110.jpg", "B111.jpg", "B112.jpg", "B113.jpg", "B114.jpg", "B115.jpg", "B116.jpg", "B117.jpg", "B118.jpg", "B119.jpg", "B12.jpg", "B120.jpg", "B13.jpg", "B14.jpg", "B15.jpg", "B16.jpg", "B17.jpg", "B18.jpg", "B19.jpg", "B2.jpg", "B20.jpg", "B21.jpg", "B22.jpg", "B23.jpg", "B24.jpg", "B25.jpg", "B26.jpg", "B27.jpg", "B101.jpg", "B29.jpg", "B3.jpg", "B30.jpg", "B31.jpg", "B32.jpg", "B33.jpg", "B34.jpg", "B35.jpg", "B36.jpg", "B37.jpg", "B38.jpg", "B39.jpg", "B4.jpg", "B40.jpg", "B41.jpg", "B42.jpg", "B43.jpg", "B44.jpg", "B45.jpg", "B46.jpg", "B47.jpg", "B48.jpg", "B49.jpg", "B5.jpg", "B50.jpg", "B51.jpg", "B52.jpg", "B53.jpg", "B54.jpg", "B55.jpg", "B56.jpg", "B57.jpg", "B58.jpg", "B59.jpg", "B6.jpg", "B60.jpg", "B61.jpg", "B62.jpg", "B63.jpg", "B64.jpg", "B65.jpg", "B66.jpg", "B67.jpg", "B68.jpg", "B69.jpg", "B7.jpg", "B70.jpg", "B71.jpg", "B72.jpg", "B73.jpg", "B74.jpg", "B75.jpg", "B76.jpg", "B77.jpg", "B78.jpg", "B79.jpg", "B8.jpg", "B80.jpg", "B81.jpg", "B82.jpg", "B83.jpg", "B84.jpg", "B85.jpg", "B86.jpg", "B87.jpg", "B28.jpg", "B89.jpg", "B9.jpg", "B90.jpg", "B91.jpg", "B92.jpg", "B93.jpg", "B94.jpg", "B95.jpg", "B96.jpg", "B97.jpg", "B98.jpg", "B99.jpg", "А1.jpg", "А10.jpg", "B88.jpg", "А101.jpg", "А102.jpg", "А11.jpg", "А12.jpg", "А13.jpg", "А14.jpg", "А15.jpg", "А16.jpg", "А17.jpg", "А18.jpg", "А19.jpg", "А2.jpg", "А20.jpg", "А21.jpg", "А22.jpg", "А23.jpg", "А24.jpg", "А25.jpg", "А26.jpg", "А27.jpg", "А28.jpg", "А29.jpg", "А3.jpg", "А30.jpg", "А31.jpg", "А32.jpg", "А33.jpg", "А34.jpg", "А35.jpg", "А36.jpg", "А37.jpg", "А38.jpg", "А39.jpg", "А4.jpg", "А40.jpg", "А41.jpg", "А42.jpg", "А43.jpg", "А44.jpg", "А45.jpg", "А46.jpg", "А47.jpg", "А48.jpg", "А49.jpg", "А5.jpg", "А50.jpg", "А51.jpg", "А52.jpg", "А53.jpg", "А54.jpg", "А55.jpg", "А56.jpg", "А57.jpg", "А58.jpg", "А59.jpg", "А6.jpg", "А60.jpg", "А61.jpg", "А62.jpg", "А63.jpg", "А64.jpg", "А65.jpg", "А66.jpg", "А67.jpg", "А68.jpg", "А69.jpg", "А7.jpg", "А70.jpg", "А71.jpg", "А72.jpg", "А73.jpg", "А74.jpg", "А75.jpg", "А76.jpg", "А77.jpg", "А78.jpg", "А79.jpg", "А8.jpg", "А100.jpg", "А80.jpg", "А81.jpg", "А82.jpg", "А84.jpg", "А85.jpg", "А86.jpg", "А87.jpg", "А88.jpg", "А89.jpg", "А9.jpg", "А90.jpg", "А91.jpg", "А92.jpg", "А93.jpg", "А94.jpg", "А95.jpg", "А96.jpg", "А97.jpg", "А98.jpg", "А99.jpg", "С1.jpg", "С10.jpg", "B10.jpg", "С101.jpg", "С102.jpg", "С103.jpg", "С11.jpg", "С12.jpg", "С13.jpg", "С14.jpg", "С15.jpg", "С16.jpg", "B100.jpg", "С18.jpg", "С19.jpg", "С2.jpg", "С20.jpg", "С21.jpg", "С24.jpg", "С25.jpg", "С26.jpg", "С27.jpg", "С28.jpg", "С29.jpg", "С3.jpg", "С30.jpg", "С32.jpg", "С34.jpg", "С35.jpg", "С36.jpg", "С37.jpg", "С17.jpg", "С39.jpg", "С4.jpg", "С40.jpg", "С41.jpg", "С42.jpg", "С43.jpg", "С44.jpg", "С45.jpg", "С46.jpg", "С47.jpg", "С48.jpg", "С49.jpg", "С5.jpg", "С51.jpg", "С53.jpg", "С54.jpg", "С55.jpg", "С56.jpg", "С57.jpg", "С58.jpg", "С59.jpg", "С6.jpg", "С60.jpg", "С61.jpg", "С62.jpg", "С63.jpg", "С64.jpg", "С65.jpg", "С66.jpg", "С67.jpg", "С69.jpg", "С7.jpg", "С70.jpg", "С73.jpg", "С74.jpg", "С75.jpg", "С38.jpg", "С100.jpg", "С78.jpg", "С79.jpg", "С8.jpg", "С80.jpg", "С76.jpg", "С82.jpg", "С83.jpg", "С84.jpg", "С85.jpg", "С86.jpg", "С88.jpg", "А83.jpg", "С9.jpg", "С90.jpg", "С91.jpg", "С92.jpg", "С93.jpg", "С94.jpg", "С95.jpg", "С96.jpg", "С89.jpg", "С98.jpg", "С99.jpg", "С97.jpg", "С81.jpg"}
	// srcDir := "/home/vlasov/folder/mgs2/dogfound/core/my_data/new_images/"
	// unparsedDdt := "/home/vlasov/folder/mgs2/dogfound/core/my_data/unparsed/"
	// if err := cv.MoveSelectedToFile(srcDir, unparsedDdt, unparsed); err != nil {
	// 	panic(err)
	// }
	// os.Exit(0)

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

	go api.ConnectorListen(":6000", imageSourceDirectory)
	api.Serve(classificatorAddress)
}
