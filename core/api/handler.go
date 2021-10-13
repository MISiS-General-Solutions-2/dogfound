package api

import (
	"net/http"
	"pet-track/database"

	"github.com/gin-gonic/gin"
)

func getImage(ctx *gin.Context) {
	name := ctx.Param("name")
	ctx.File(database.GetFilePath(name))
}
func getImagesByFeatures(ctx *gin.Context) {
	var features map[string]interface{}
	if err := ctx.BindJSON(&features); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	imgs, err := database.GetImagesByFeatures(features)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, imgs)
}
