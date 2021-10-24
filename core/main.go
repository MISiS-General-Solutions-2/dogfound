package main

import (
	"dogfound/api"
	"dogfound/database"
	"dogfound/processor"
	"log"
)

func main() {
	go func() {
		closer := database.Connect()
		defer closer()
		for {
			if err := processor.ProcessNewImages(); err != nil {
				log.Printf("error occurred processing images: %v\nRetrying...", err)
			}
		}
	}()

	api.Serve()

}
