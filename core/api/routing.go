package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	classificatorAddress string
	volunteerFolder      string
)

func Serve(classificator, volunteerImageFolder string) {
	classificatorAddress = classificator
	volunteerFolder = volunteerImageFolder

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Access-Control-Allow-Headers", "Authorization", "X-Requested-With"},
		MaxAge:       12 * time.Hour,
	}))

	api := router.Group("/api")
	api.POST("/image/by-classes", getImagesByClasses)
	api.GET("/image/:name",
		getImage)
	api.POST("/image/similar",
		getSimilar)
	api.PUT("/image/upload", upload)

	router.Run(":5000")
}
