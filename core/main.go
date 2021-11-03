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
	simpleParseArgs()
	fmt.Println(sampleInterval)

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
