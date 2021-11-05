package api

import (
	"dogfound/cv"
	"dogfound/database"
	doghttp "dogfound/http"
	"image/color"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func getImage(ctx *gin.Context) {
	name := ctx.Param("name")
	omit, _ := strconv.Atoi(ctx.Query("omit_crop"))
	a, err := database.GetAdditionalInfo(name)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	image := database.GetImagePath(name)
	if omit != 1 {
		imgWithCrop, err := cv.DrawRect(image, a.Crop, color.RGBA{0, 255, 0, 128}, 2)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.Data(http.StatusOK, "image/jpeg", imgWithCrop)
	} else {
		ctx.File(image)
	}
}

func getImagesByClasses(ctx *gin.Context) {
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
func getSimilar(ctx *gin.Context) {
	t0, _ := strconv.Atoi(ctx.Query("t0"))
	t1, _ := strconv.Atoi(ctx.Query("t1"))

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	files := form.File["file"]
	if len(files) != 1 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "must provide exactly one image")
		return
	}

	tempDir, err := os.MkdirTemp("", "dogfound")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	filename := tempDir + files[0].Filename
	if err = ctx.SaveUploadedFile(files[0], filename); err != nil {
		os.Remove(tempDir)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer func() {
		os.Remove(filename)
		os.Remove(tempDir)
	}()

	cr, err := doghttp.Categorize(doghttp.Destination{Address: classificatorAddress, Retries: 3, RetryInterval: 1 * time.Second}, filename)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if cr.IsAnimal == 0 {
		ctx.JSON(http.StatusOK, SimilarResponse{IsAnimal: 0})
		return
	}
	sr, err := database.GetImagesByClasses(dbRequestFromCategories(cr, t0, t1))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, SimilarResponse{IsAnimal: 1, Results: sr})
}
func dbRequestFromCategories(ct doghttp.CategorizationResponse, t0, t1 int) map[string]interface{} {
	res := make(map[string]interface{})
	res["tail"] = float64(ct.Tail)
	res["color"] = float64(ct.Color)
	if t0 != 0 {
		res["t0"] = float64(t0)
	}
	if t1 != 0 {
		res["t1"] = float64(t1)
	}
	return res
}
