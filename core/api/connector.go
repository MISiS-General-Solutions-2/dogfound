package api

import (
	"dogfound/shared"
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ConnectorListen(addr, dst string) {
	router := gin.Default()

	router.PUT("/image/upload", func(ctx *gin.Context) {
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		files := form.File["file"]
		if err := checkAndTryFixExtensions(files); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}
		for _, file := range files {
			ext := shared.GetExtension(file.Filename)
			filename := uuid.NewString() + ext
			if err := ctx.SaveUploadedFile(file, dst+filename); err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
		ctx.Status(http.StatusNoContent)
	})

	router.Run(addr)
}
func checkAndTryFixExtensions(files []*multipart.FileHeader) error {
	for _, file := range files {
		switch shared.GetExtension(file.Filename) {
		case ".jpeg":
			file.Filename = shared.ChangeExtension(file.Filename, ".jpg")
		case ".jpg":
		default:
			return errors.New("only .jpg format allowed")
		}
	}
	return nil
}
