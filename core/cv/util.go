package cv

import (
	"log"
	"os"
	"strings"
)

func MoveSelectedToFile(srcDir, dstDir string, imgs []string) error {
	for _, img := range imgs {
		if err := os.Rename(srcDir+img, dstDir+img); err != nil {
			return err
		}
	}
	return nil
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
