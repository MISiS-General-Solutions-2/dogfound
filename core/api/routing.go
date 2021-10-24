package api

import (
	"github.com/gin-gonic/gin"
)

func Serve() {
	router := gin.Default()

	api := router.Group("/api")
	api.POST("/image/by-classes", getImagesByFeatures)
	api.GET("/image/:name", getImage)

	router.Run(":5000")
}
