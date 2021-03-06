package cv

import (
	"errors"
	"fmt"
	"image"

	"github.com/otiai10/gosseract"
	"gocv.io/x/gocv"
)

var imshow bool

func ParseImage(img string) (camID string, timestamp int64, err error) {
	//imshow = true
	// if img != "/opt/dogfound/data/new_images/B37.jpg" {
	// 	return
	// }
	functions := []func(img []byte){
		func(img []byte) {
			camID = tessParseCamIDFromBlackTop(img)
		},
		func(img []byte) {
			timestamp = parseTimestamp(img)
		},
	}
	// it is assumed number always has same position and pixel sizes
	rois := []image.Rectangle{
		{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 240, Y: 18}},
		{Min: image.Point{X: 0, Y: 47}, Max: image.Point{X: 240, Y: 64}},
	}
	imgMat := gocv.IMRead(img, gocv.IMReadGrayScale)
	defer imgMat.Close()

	for i := 0; i < 2; i++ {
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
	if timestamp == 0 {
		//timestamp = parseTimestampTopRight(imgMat)
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
func isBlackHeader(img gocv.Mat, blackValueThresh, whiteValueThresh int, blackPartThresh, whitePartThresh float64) bool {
	blackCount := 0
	whiteCount := 0
	for y := 0; y < img.Rows(); y++ {
		for x := 0; x < img.Cols(); x++ {
			if int(img.GetUCharAt(y, x)) <= blackValueThresh {
				blackCount += 1
			} else if int(img.GetUCharAt(y, x)) >= whiteValueThresh {
				whiteCount += 1
			}
		}
	}
	total := float64(img.Cols() * img.Rows())
	blackPercent := float64(blackCount) / total
	whitePercent := float64(whiteCount) / total

	return blackPercent >= blackPartThresh && whitePercent >= whitePartThresh
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
	if imshow {
		w := gocv.NewWindow("cropped")
		defer w.Close()
		w.IMShow(*cropped)
		w.WaitKey(0)
	}
	if !isBlackHeader(*cropped, 5, 250, 0.5, 0) {
		return nil, nil
	}
	processed := preProcess(*cropped)
	if imshow {
		w := gocv.NewWindow("processed")
		defer w.Close()
		w.IMShow(processed)
		w.WaitKey(0)
	}
	buf, err := gocv.IMEncode(gocv.FileExt(format), processed)
	if err != nil {
		return nil, err
	}
	return buf.GetBytes(), nil
}
func omitBlackHeader(img gocv.Mat) *gocv.Mat {
	for i := 1; i < img.Rows(); i++ {
		crop := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: img.Cols() - 1, Y: i}}
		cropped := getCroppedPart(img, crop)
		if cropped == nil {
			return nil
		}

		if !isBlackHeader(*cropped, 5, 250, 0.9, 0) {
			start := i
			if start-5 > 0 {
				start = i - 5
			}
			crop := image.Rectangle{Min: image.Point{X: 0, Y: start}, Max: image.Point{X: img.Cols() - 1, Y: img.Rows() - 1}}
			return getCroppedPart(img, crop)
		}
	}
	return nil
}
func selectTopPart(img gocv.Mat, part float64) *gocv.Mat {
	crop := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: img.Cols() - 1, Y: int(float64(img.Cols()) * part)}}
	return getCroppedPart(img, crop)
}
func preProcessTopRight(img gocv.Mat) gocv.Mat {
	w := gocv.NewWindow("processing")
	defer w.Close()

	gocv.Resize(img, &img, image.Point{}, 3, 3, gocv.InterpolationCubic)

	// morhtElem := gocv.GetStructuringElement(gocv.MorphShape(gocv.MorphEllipse), image.Point{2, 2})
	// defer morhtElem.Close()
	// gocv.MorphologyEx(img, &img, gocv.MorphDilate, morhtElem)

	gocv.MedianBlur(img, &img, 5)

	gocv.Threshold(img, &img, 200, 255, gocv.ThresholdBinaryInv)
	w.IMShow(img)
	w.WaitKey(0)
	return img
}
func parseTimestampTopRight(img gocv.Mat) int64 {
	headerLess := omitBlackHeader(img)
	if headerLess == nil {
		fmt.Println("could not omit black header")
		return 0
	}

	roi := selectTopPart(*headerLess, 0.09)
	if roi == nil {
		fmt.Println("could not sect top part")
		return 0
	}

	processed := preProcessTopRight(*roi)

	buf, err := gocv.IMEncode(".png", processed)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	b := buf.GetBytes()
	if len(b) == 0 {
		return 0
	}

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImageFromBytes(b)
	client.SetWhitelist("0123456789--")
	text, err := client.Text()
	if err := client.SetVariable("oem", "2"); err != nil {
		fmt.Println(err)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println(text)

	return 0
}
