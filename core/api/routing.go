package api

import (
	"github.com/gin-gonic/gin"
)

func Serve() {
	router := gin.Default()

	api := router.Group("/api")
	api.GET("/image/by-features", getImagesByFeatures)
	api.GET("/image/:name", getImage)

	router.Run(":6000")
}
