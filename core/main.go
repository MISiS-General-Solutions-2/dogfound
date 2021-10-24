package main

import (
	"log"
	"pet-track/api"
	"pet-track/database"
	"pet-track/processor"
)

// @title Blueprint Swagger API
// @version 1.0
// @description Swagger API for Golang Project Blueprint.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email martin7.heinz@gmail.com

// @license.name MIT
// @license.url https://github.com/MartinHeinz/go-project-blueprint/blob/master/LICENSE

// @BasePath /api/v1
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
