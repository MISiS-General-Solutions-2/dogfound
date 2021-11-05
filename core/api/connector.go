package api

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
		if err := checkExtensions(files); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}
		for _, file := range files {
			if err := ctx.SaveUploadedFile(file, dst+file.Filename); err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
		ctx.Status(http.StatusNoContent)
	})

	router.Run(addr)
}
func checkExtensions(files []*multipart.FileHeader) error {
	for _, file := range files {
		if idx := strings.LastIndexByte(file.Filename, '.'); idx != -1 {
			switch file.Filename[idx:] {
			case ".jpeg":
				file.Filename = file.Filename[:idx] + ".jpg"
			case ".jpg":
			default:
				return errors.New("only .jpg format allowed")
			}
		}
	}
	return nil
}
