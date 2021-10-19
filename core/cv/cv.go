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
	"os"
	"pet-track/database"
	"regexp"
	"strings"
	"unsafe"

	"gocv.io/x/gocv"
)

const (
	gocr = `D:\Papka\myProgramms\gocr\gocr049.exe`
)

var (
	reCameraName = regexp.MustCompile(`[A-Z]+?_\S+?_\S+`)

	apxi8 = gocv.IMRead(database.DataPath+"apxi8", gocv.IMReadGrayScale)
)

func CParseImage(img []byte, args []string) string {
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

func GetImagesText(imgs []string) ([]string, error) {
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
		if err := RetrieveBlackTop(img, cb); err != nil {
			return nil, err
		}
	}

	return parsed, nil
}
func parseSingleImageDirect(img []byte) (string, error) {
	out := CParseImage(img, []string{"gocv", "-C", "0-9a-zA-Z"})
	return ParseRecognized(strings.NewReader(out)), nil
}

// func parseSingleImage(img []byte) (string, error) {
// 	if err := os.WriteFile(`database.TempDir`+`top.pgm`, img, 0600); err != nil {
// 		return "", err
// 	}
// 	defer os.Remove(database.TempDir + `top.pgm`)

// 	cmd := exec.Command(gocr, "-i", database.TempDir+`top.pgm`, "-C", "0-9a-zA-Z")
// 	out, err := cmd.StdoutPipe()
// 	if err != nil {
// 		return "", err
// 	}
// 	if err = cmd.Start(); err != nil {
// 		log.Println(err)
// 	}
// 	return ParseRecognized(out), nil
// }
func ParseRecognized(r io.Reader) string {
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
func RetrieveBlackTop(file string, cb func([]byte) bool) error {

	img := gocv.IMRead(file, gocv.IMReadGrayScale)
	defer img.Close()

	cropRect := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 240, Y: 18}}

	// drop too small images
	if img.Cols() < cropRect.Max.X || img.Rows() < cropRect.Max.Y {
		cb(nil)
		//MoveToFile(file, database.data `/special/`)
		return nil
	}

	// select region of interest
	cropped := img.Region(cropRect)
	thresh := 5
	// drop images with not black header
	shouldBeBlack := cropped.ColRange(0, 5)
	if int(shouldBeBlack.Mean().Val1) > thresh {
		log.Println("not black")
		cb(nil)
		//MoveToFile(file, `../data/special/`)
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

	if strings.Contains(file, target) {
		window.IMShow(cropped)
		gocv.WaitKey(0)
	}

	morhtElem := gocv.GetStructuringElement(gocv.MorphShape(gocv.MorphEllipse), image.Point{2, 2})
	defer morhtElem.Close()
	gocv.MorphologyEx(cropped, &cropped, gocv.MorphDilate, morhtElem)

	if strings.Contains(file, target) {
		window.IMShow(cropped)
		gocv.WaitKey(0)
	}

	gocv.GaussianBlur(cropped, &cropped, image.Point{5, 5}, 0, 0, gocv.BorderDefault)

	// window.IMShow(cropped)
	// window.WaitKey(0)

	gocv.Threshold(cropped, &cropped, 128, 255, gocv.ThresholdBinaryInv)

	if strings.Contains(file, target) {
		window.IMShow(cropped)
		gocv.WaitKey(0)
	}

	// window.IMShow(cropped)
	// window.WaitKey(0)

	buf, err := gocv.IMEncode(".pgm", cropped)
	if err != nil {
		return err
	}
	cb(buf.GetBytes())
	return nil
	// found := cb(buf.GetBytes())

	// if !found {
	// 	//MoveToFile(file, `../data/special/`)
	// }
	// return nil
	// if err != nil {
	// 	return err
	// }
	// if len(buf.GetBytes()) == 0 {
	// 	fmt.Println(file)
	// }
	// found := cb(buf.GetBytes())
	// var top gocv.Mat
	// {
	// 	i := 0
	// 	right := img.ColRange(img.Cols()-10, img.Cols())
	// 	for j := 0; j < img.Rows(); j++ {
	// 		rowSample := right.RowRange(j, j+1)
	// 		if int(rowSample.Mean().Val1) < thresh {
	// 			i++
	// 		} else {
	// 			break
	// 		}
	// 	}
	// 	// contains some black top
	// 	if i < 5 {
	// 		log.Printf("non standard format for file %v\n", file)
	// 		// err := os.Rename(file, `../data/special`+file[strings.LastIndexByte(file, '\\')+1:])
	// 		// if err != nil {
	// 		// 	log.Fatal(err)
	// 		// }
	// 		cb(nil)
	// 		return nil
	// 	}
	// 	top = img.RowRange(0, i/4)
	// }
	// defer top.Close()

	// buf, err := gocv.IMEncode(".pgm", top)
	// if err != nil {
	// 	return err
	// }
	// if len(buf.GetBytes()) == 0 {
	// 	fmt.Println(file)
	// }
	// found := cb(buf.GetBytes())
	// if !found {
	// 	// window := gocv.NewWindow("unparsed")
	// 	// defer window.Close()

	// 	// window.IMShow(top)
	// 	// window.WaitKey(0)
	// 	var filename string
	// 	if idx := strings.LastIndexByte(file, '\\'); idx != -1 {
	// 		filename = file[idx+1:]
	// 	} else if idx := strings.LastIndexByte(file, '/'); idx != -1 {
	// 		filename = file[idx+1:]
	// 	}
	// 	err := os.Rename(file, `../data/special/`+filename)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// return nil
}
func MoveToFile(file, to string) {
	var filename string
	if idx := strings.LastIndexByte(file, '\\'); idx != -1 {
		filename = file[idx+1:]
	} else if idx := strings.LastIndexByte(file, '/'); idx != -1 {
		filename = file[idx+1:]
	}
	err := os.Rename(file, to+filename)
	if err != nil {
		log.Fatal(err)
	}
}
func RetrieveAllFromSpecial() {
	de, err := os.ReadDir(`../data/special/`)
	if err != nil {
		panic(err)
	}
	for _, e := range de {
		MoveToFile(`../data/special/`+e.Name(), `../data/img/`)
	}
}
