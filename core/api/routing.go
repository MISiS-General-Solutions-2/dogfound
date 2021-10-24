package api

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Serve() {
	router := gin.Default()

	api := router.Group("/api")
	api.GET("/image/by-classes", getImagesByFeatures)
	api.GET("/image/:name", getImage)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":6000")
}
