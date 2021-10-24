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
	"dogfound/database"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"unsafe"

	"gocv.io/x/gocv"
)

var (
	reCameraName = regexp.MustCompile(`[A-Z]+?_\S+?_\S+`)

	apxi8 = gocv.IMRead(database.DataPath+"apxi8", gocv.IMReadGrayScale)
)

func cParseCamID(img []byte, args []string) string {
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
func parseCamID(img []byte) (string, error) {
	out := cParseCamID(img, []string{"gocv", "-C", "0-9a-zA-Z"})
	return parseRecognizedCamID(strings.NewReader(out)), nil
}
func parseRecognizedCamID(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return ""
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
