package cv

/*
#cgo LDFLAGS: -L${SRCDIR}/../libs -lPgm2asc
#include <stdlib.h>
char *parse_pgm(int img_data_len, char *image, int argn, char *argv[]);
*/
import "C"
import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"
	"pet-track/database"
	"regexp"
	"strings"
	"unsafe"

	"gocv.io/x/gocv"
)

var (
	reCameraName = regexp.MustCompile(`[A-Z]+?_\S+?_\S+`)

	apxi8 = gocv.IMRead(database.DataPath+"apxi8", gocv.IMReadGrayScale)
)

func cParseImage(img []byte, args []string) string {
	argv := make([]*C.char, len(args))
	for i, s := range args {
		cs := C.CString(s)
		defer C.free(unsafe.Pointer(cs))
		argv[i] = cs
	}

	cres := C.parse_pgm(C.int(len(img)), (*C.char)(unsafe.Pointer(&img[0])), C.int(len(args)), &argv[0])
	res := C.GoString(cres)
	C.free(unsafe.Pointer(cres))
	return res
}

func GetImagesCamIDs(directory string, imgs []string) ([]string, error) {
	parsed := make([]string, 0, len(imgs))
	cb := func(top []byte) bool {
		if len(top) == 0 {
			parsed = append(parsed, "")
			return false
		}
		s, err := parseSingleImageDirect(top)
		if err != nil {
			panic(err)
		}
		parsed = append(parsed, s)
		return s != ""
	}
	for _, img := range imgs {
		if err := retrieveBlackTop(directory+img, cb); err != nil {
			return nil, err
		}
	}

	return parsed, nil
}
func parseSingleImageDirect(img []byte) (string, error) {
	out := cParseImage(img, []string{"gocv", "-C", "0-9a-zA-Z"})
	return parseRecognized(strings.NewReader(out)), nil
}

func parseRecognized(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	n := len(b)

	scan := bufio.NewScanner(bytes.NewReader(b))
	for {
		line := scan.Bytes()
		if match := reCameraName.Find(line); match != nil {
			return string(fixCameraName(match))
		}
		if !scan.Scan() {
			break
		}
	}
	if n != 0 {
		fmt.Printf("%s\n", b)
	}
	return ""
}
func fixCameraName(b []byte) []byte {
	return fixLowCase(fixNumbers(fixStart(b)))
}
func fixLowCase(b []byte) []byte {
	upper := strings.ToUpper(string(b))
	return []byte(strings.ReplaceAll(upper, "HD", "hd"))
}
func fixStart(b []byte) []byte {
	if bytes.HasPrefix(b, []byte("N_")) {
		res := []byte("PV")
		return append(res, b...)
	}
	if bytes.HasPrefix(b, []byte("VN_")) {
		res := []byte("P")
		return append(res, b...)
	}
	if bytes.HasPrefix(b, []byte("D_N_")) {
		res := []byte("DVN")
		return append(res, b[3:]...)
	}
	return b
}
func fixNumbers(b []byte) []byte {
	idx := 0
	if bytes.Contains(b, []byte("hd")) {
		for i := 0; i < 3; i++ {
			idx = bytes.IndexByte(b[idx:], '_') + idx + 1
			if idx == 0 {
				return b
			}
		}
	} else {
		for i := 0; i < 2; i++ {
			idx = bytes.IndexByte(b[idx:], '_') + idx + 1
			if idx == 0 {
				return b
			}
		}
	}
	s := string(b[idx:])
	s = strings.ReplaceAll(s, "o", "0")
	s = strings.ReplaceAll(s, "O", "0")
	s = strings.ReplaceAll(s, "D", "0")
	s = strings.ReplaceAll(s, "l", "1")
	s = strings.ReplaceAll(s, "L", "1")
	s = strings.ReplaceAll(s, "q", "4")
	s = strings.ReplaceAll(s, "B", "8")

	return append(b[:idx], []byte(s)...)
}

// it is assumed number always has same position and pixel sizes
func retrieveBlackTop(file string, cb func([]byte) bool) error {
	img := gocv.IMRead(file, gocv.IMReadGrayScale)
	defer img.Close()

	cropRect := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 240, Y: 18}}

	// drop too small images
	if img.Cols() < cropRect.Max.X || img.Rows() < cropRect.Max.Y {
		cb(nil)
		//MoveToFile(file, `data/unparsed/`)
		return nil
	}

	// select region of interest
	cropped := img.Region(cropRect)
	thresh := 10
	// drop images with not black header
	shouldBeBlack := cropped.ColRange(0, 5)
	if int(shouldBeBlack.Mean().Val1) > thresh {
		log.Printf("not black with thresh %v\n", int(shouldBeBlack.Mean().Val1))
		cb(nil)
		//MoveToFile(file, `data/unparsed/`)
		return nil
	}

	// preprocessing
	var window *gocv.Window
	target := "none"
	if strings.Contains(file, target) {
		window = gocv.NewWindow(file)
	}
	// apxib helps somehow
	gocv.Hconcat(apxi8, cropped, &cropped)

	if strings.Contains(file, target) {
		window.IMShow(cropped)
		gocv.WaitKey(0)
	}

	gocv.Resize(cropped, &cropped, image.Point{}, 3, 3, gocv.InterpolationCubic)

	morhtElem := gocv.GetStructuringElement(gocv.MorphShape(gocv.MorphEllipse), image.Point{2, 2})
	defer morhtElem.Close()
	gocv.MorphologyEx(cropped, &cropped, gocv.MorphDilate, morhtElem)

	if strings.Contains(file, target) {
		window.IMShow(cropped)
		gocv.WaitKey(0)
	}

	gocv.GaussianBlur(cropped, &cropped, image.Point{5, 5}, 0, 0, gocv.BorderDefault)

	gocv.Threshold(cropped, &cropped, 128, 255, gocv.ThresholdBinaryInv)

	buf, err := gocv.IMEncode(".pgm", cropped)
	if err != nil {
		return err
	}
	cb(buf.GetBytes())
	//MoveToFile(file, `data/unparsed/`)
	return nil
}
