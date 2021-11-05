package cv

import (
	"dogfound/shared"
	"errors"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

func DrawRect(img string, rect [4]int, col color.RGBA, thickness int) ([]byte, error) {
	ext := shared.GetExtension(img)
	if ext == "" {
		return nil, errors.New("file has no extension")
	}
	imgMat := gocv.IMRead(img, gocv.IMReadColor)
	if imgMat.Empty() {
		return nil, errors.New("empty image")
	}
	gocv.Rectangle(&imgMat, image.Rect(rect[0], rect[1], rect[2], rect[3]), col, thickness)

	buf, err := gocv.IMEncode(gocv.FileExt(ext), imgMat)
	if err != nil {
		return nil, err
	}
	return buf.GetBytes(), nil
}
