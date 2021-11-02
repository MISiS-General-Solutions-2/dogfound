package api

import (
	"dogfound/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getImage(ctx *gin.Context) {
	name := ctx.Param("name")
	ctx.File(database.GetImagePath(name))
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
