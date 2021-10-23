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
func initDB() {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS images
	(
		filename TEXT NOT NULL PRIMARY KEY,
		is_animal_there INTEGER,
		is_it_a_dog INTEGER,
		is_the_owner_there INTEGER,
		color INTEGER,
		tail INTEGER,
		cam_id TEXT,
		timestamp INTEGER
	);
	CREATE TABLE IF NOT EXISTS registries
	(
		cam_id TEXT NOT NULL PRIMARY KEY,
		address TEXT NOT NULL,
		lat REAL NOT NULL,
		lon REAL NOT NULL
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		db.Close()
		panic(err)
	}
	if err = populateRegistries(); err != nil {
		panic(err)
	}
}
