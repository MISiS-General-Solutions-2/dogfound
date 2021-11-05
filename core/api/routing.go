package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var classificatorAddress string

func Serve(classificator string) {
	classificatorAddress = classificator

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

	router.Run(":5000")
}
