package cv

import (
	"errors"
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

func ParseImage(img string) (camID string, timestamp int64, err error) {
	functions := []func(img []byte){
		func(img []byte) {
			camID = tessParseCamID(img)
		},
		func(img []byte) {
			timestamp = parseTimestamp(img)
		},
	}
	// it is assumed number always has same position and pixel sizes
	rois := []image.Rectangle{
		{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 240, Y: 18}},
		{Min: image.Point{X: 0, Y: 53}, Max: image.Point{X: 240, Y: 64}},
	}
	imgMat := gocv.IMRead(img, gocv.IMReadGrayScale)
	defer imgMat.Close()

	for i := range functions {
		var imgBytes []byte
		imgBytes, err = getProcessedRegionAsBytes(imgMat, rois[i], ".png")
		if err != nil {
			fmt.Println(img)
			return
		}
		if len(imgBytes) == 0 {
			continue
		}

		functions[i](imgBytes)
	}
	return
}
func getCroppedPart(img gocv.Mat, crop image.Rectangle) *gocv.Mat {
	// drop too small images
	if img.Cols() < crop.Max.X || img.Rows() < crop.Max.Y {
		return nil
	}

	// select region of interest
	cropped := img.Region(crop)
	return &cropped
}
func isRegionMedianBelowThresh(img gocv.Mat, thresh int) bool {
	return int(img.Mean().Val1) <= thresh
}
func preProcess(img gocv.Mat) gocv.Mat {
	gocv.Resize(img, &img, image.Point{}, 3, 3, gocv.InterpolationCubic)

	morhtElem := gocv.GetStructuringElement(gocv.MorphShape(gocv.MorphEllipse), image.Point{2, 2})
	defer morhtElem.Close()
	gocv.MorphologyEx(img, &img, gocv.MorphDilate, morhtElem)

	gocv.GaussianBlur(img, &img, image.Point{5, 5}, 0, 0, gocv.BorderDefault)

	gocv.Threshold(img, &img, 128, 255, gocv.ThresholdBinaryInv)

	return img
}

func getProcessedRegionAsBytes(img gocv.Mat, crop image.Rectangle, format string) ([]byte, error) {

	cropped := getCroppedPart(img, crop)
	if cropped == nil {
		return nil, errors.New("could not crop image")
	}
	if !isRegionMedianBelowThresh(*cropped, 20) {
		return nil, nil
	}
	processed := preProcess(*cropped)
	buf, err := gocv.IMEncode(gocv.FileExt(format), processed)
	if err != nil {
		return nil, err
	}
	return buf.GetBytes(), nil
}
