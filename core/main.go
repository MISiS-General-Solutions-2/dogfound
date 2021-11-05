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
	cameraSourceDirectory    = database.DataPath + "new_images/"
	volunteerSourceDirectory = database.DataPath + "volunteer_added/"
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
	fmt.Println(classificatorAddress)

	close := database.Connect()
	defer close()

	processor.Processor = *processor.CreateProcessor(&processor.Config{
		Classificator:           http.Destination{Address: classificatorAddress, Retries: 3, RetryInterval: 1 * time.Second},
		CameraInputDirectory:    cameraSourceDirectory,
		VolunteerInputDirectory: volunteerSourceDirectory,
		NumWorkers:              numWorkers,
		SampleInterval:          time.Duration(sampleInterval) * time.Second,
	})
	go processor.Processor.Start()

	go api.ConnectorListen(":6000", cameraSourceDirectory)
	api.Serve(classificatorAddress, volunteerSourceDirectory)
}
