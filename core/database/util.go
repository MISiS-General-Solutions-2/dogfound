package database

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
)

//utility to flaten dataset
func PopulateWithImages(path string) error {
	reFileName := regexp.MustCompile(`^\d+\.jpg`)
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !reFileName.MatchString(info.Name()) {
				return nil
			}

			dest, err := os.OpenFile(imagePath+info.Name(), os.O_CREATE|os.O_RDWR, 0600)
			if err != nil {
				return err
			}
			defer dest.Close()

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(dest, file)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}
