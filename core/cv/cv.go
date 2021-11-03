package cv

import (
	"fmt"
	"image"
	"log"

	"gocv.io/x/gocv"
)

func ParseImage(directory string, img string) (camID string, timestamp int64, err error) {
	camIDCb := func(img []byte) (bool, error) {
		if len(img) == 0 {
			camID = ""
			return false, nil
		}
		fmt.Println(len(img))
		s, err := parseCamID(img)
		if err != nil {
			return false, err
		}
		camID = s
		return s != "", nil
	}
	timestampCb := func(img []byte) (bool, error) {
		if len(img) == 0 {
			timestamp = 0
			return false, nil
		}
		s, err := parseTimestamp(img)
		if err != nil {
			return false, err
		}
		timestamp = s
		return s != 0, nil
	}
	if err = retrieveBlackTop(directory+img, camIDCb, timestampCb); err != nil {
		return
	}
	return
}
func getCroppedPart(img *gocv.Mat, crop image.Rectangle) *gocv.Mat {
	// drop too small images
	if img.Cols() < crop.Max.X || img.Rows() < crop.Max.Y {
		return nil
	}

	// select region of interest
	cropped := img.Region(crop)
	return &cropped
}
func isRegionMedianBelowThresh(img *gocv.Mat, thresh int) bool {
	shouldBeBlack := img.ColRange(0, 5)
	return int(shouldBeBlack.Mean().Val1) <= thresh
}

func cropAndPreProcess(img *gocv.Mat, crop image.Rectangle, addApxib bool) *gocv.Mat {
	// drop too small images
	if img.Cols() < crop.Max.X || img.Rows() < crop.Max.Y {
		return nil
	}

	// select region of interest
	cropped := img.Region(crop)
	thresh := 10
	// drop images with not black header
	shouldBeBlack := cropped.ColRange(0, 5)
	if int(shouldBeBlack.Mean().Val1) > thresh {
		log.Printf("not black with thresh %v\n", int(shouldBeBlack.Mean().Val1))
		return nil
	}

	// preprocessing
	if addApxib {
		// apxib helps somehow
		gocv.Hconcat(apxi8, cropped, &cropped)
	}

	gocv.Resize(cropped, &cropped, image.Point{}, 3, 3, gocv.InterpolationCubic)

	morhtElem := gocv.GetStructuringElement(gocv.MorphShape(gocv.MorphEllipse), image.Point{2, 2})
	defer morhtElem.Close()
	gocv.MorphologyEx(cropped, &cropped, gocv.MorphDilate, morhtElem)

	gocv.GaussianBlur(cropped, &cropped, image.Point{5, 5}, 0, 0, gocv.BorderDefault)

	gocv.Threshold(cropped, &cropped, 128, 255, gocv.ThresholdBinaryInv)

	return &cropped
}

// it is assumed number always has same position and pixel sizes
func retrieveBlackTop(file string, camIDCb, timestampCb func([]byte) (bool, error)) error {
	img := gocv.IMRead(file, gocv.IMReadGrayScale)
	defer img.Close()

	camIDRect := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 240, Y: 18}}
	// timestampRect := image.Rectangle{Min: image.Point{X: 0, Y: 53}, Max: image.Point{X: 240, Y: 64}}

	camIDCrop := cropAndPreProcess(&img, camIDRect, true)
	if camIDCrop != nil {
		buf, err := gocv.IMEncode(".pgm", *camIDCrop)
		if err != nil {
			return err
		}
		if _, err := camIDCb(buf.GetBytes()); err != nil {
			return err
		}
	} else {
		camIDCb(nil)
	}

	// timestampCrop := cropAndPreProcess(&img, timestampRect, false)
	// if timestampCrop != nil {
	// 	buf, err := gocv.IMEncode(".png", *timestampCrop)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if _, err := timestampCb(buf.GetBytes()); err != nil {
	// 		return err
	// 	}
	// } else {
	// 	timestampCb(nil)
	// }

	return nil
}
