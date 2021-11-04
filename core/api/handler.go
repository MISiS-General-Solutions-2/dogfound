package api

import (
	"dogfound/cv"
	"dogfound/database"
	"image/color"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getImage(ctx *gin.Context) {
	name := ctx.Param("name")
	a, err := database.GetAdditionalInfo(name)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	imgWithCrop, err := cv.DrawRect(database.GetImagePath(name), a.Crop, color.RGBA{0, 255, 0, 128}, 2)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Data(http.StatusOK, "image/jpeg", imgWithCrop)
}

func getImagesByFeatures(ctx *gin.Context) {
	var req map[string]interface{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := database.ValidateRequest(req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	imgs, err := database.GetImagesByClasses(req)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, imgs)
}
