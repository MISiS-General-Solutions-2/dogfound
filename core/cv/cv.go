package cv

/*
#cgo LDFLAGS: -L${SRCDIR}/../../libs -lPgm2asc
#include <stdlib.h>
char *parse_pgm(int img_data_len, char *image, int argn, char *argv[]);
*/
import "C"
import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"pet-track/database"
	"regexp"
	"unsafe"

	"gocv.io/x/gocv"
)

const (
	gocr = `D:\Papka\myProgramms\gocr\gocr049.exe`
)

var (
	reCameraName = regexp.MustCompile(`[A-Z]*?_\S*?_\S*`)
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
	return CParseImage(img, []string{"-C", "0-9a-zA-Z"}), nil
}
func parseSingleImage(img []byte) (string, error) {
	if err := os.WriteFile(database.TempDir+`merged.pgm`, img, 0600); err != nil {
		return "", err
	}
	defer os.Remove(database.TempDir + `merged.pgm`)

	cmd := exec.Command(gocr, "-i", database.TempDir+`merged.pgm`, "-C", "0-9a-zA-Z")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err = cmd.Start(); err != nil {
		log.Println(err)
	}
	return ParseRecognized(out), nil
}
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
	return bytes.ReplaceAll(b, []byte{'l'}, []byte{'1'})
}
func RetrieveBlackTop(file string, cb func([]byte) bool) error {

	img := gocv.IMRead(file, gocv.IMReadGrayScale)
	defer img.Close()

	var top gocv.Mat
	thresh := 5
	{
		i := 0
		right := img.ColRange(img.Cols()-10, img.Cols())
		for j := 0; j < img.Rows(); j++ {
			rowSample := right.RowRange(j, j+1)
			if int(rowSample.Mean().Val1) < thresh {
				i++
			} else {
				break
			}
		}
		// contains some black top
		if i < 5 {
			log.Printf("non standard format for file %v\n", file)
			// err := os.Rename(file, `../data/special`+file[strings.LastIndexByte(file, '\\')+1:])
			// if err != nil {
			// 	log.Fatal(err)
			// }
			cb(nil)
			return nil
		}
		top = img.RowRange(0, i/4)
	}
	defer top.Close()

	buf, err := gocv.IMEncode(".pgm", top)
	if err != nil {
		return err
	}
	if len(buf.GetBytes()) == 0 {
		fmt.Println(file)
	}
	cb(buf.GetBytes())
	return nil
}
