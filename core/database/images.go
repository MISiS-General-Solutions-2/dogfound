package database

import (
	"errors"
	"io/ioutil"
	"os"
)

const (
	imagePath = DataPath + "img/"
)

func GetImagePath(name string) string {
	return imagePath + name
}
func getFilesInDirectory(dir string) (string, []string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", nil, err
	}
	result := make([]string, len(files))
	for i, f := range files {
		result[i] = f.Name()
	}
	return dir, result, nil
}
func GetImages() (string, []string, error) {
	return getFilesInDirectory(imagePath)
}
func imageExists(filename string) bool {
	if _, err := os.Stat(GetImagePath(filename)); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
