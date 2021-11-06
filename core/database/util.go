package database

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
)

//utility to flaten dataset
func populateWithImages(path string) error {
	reFileName := regexp.MustCompile(`^[^\.].+\.jpg`)
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
		is_animal_there INTEGER NOT NULL,
		is_it_a_dog INTEGER NOT NULL,
		is_the_owner_there INTEGER NOT NULL,
		color INTEGER NOT NULL,
		tail INTEGER NOT NULL,
		cam_id TEXT NOT NULL,
		timestamp INTEGER NOT NULL,
		breed TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS registries
	(
		cam_id TEXT NOT NULL PRIMARY KEY,
		address TEXT NOT NULL,
		lon REAL NOT NULL,
		lat REAL NOT NULL
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		db.Close()
		panic(err)
	}
	AddAdditionalDataTable()
	if err = PopulateRegistries(); err != nil {
		panic(err)
	}
}
func AddAdditionalDataTable() {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS additional
	(
		filename TEXT NOT NULL PRIMARY KEY,
		crop_x0 INTEGER,
		crop_y0 INTEGER,
		crop_x1 INTEGER,
		crop_y1 INTEGER
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		db.Close()
		panic(err)
	}
}
func AddVolunteerSourced() {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS volunteer_sourced
	(
		filename TEXT NOT NULL PRIMARY KEY,
		lon REAL NOT NULL,
		lat REAL NOT NULL
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		db.Close()
		panic(err)
	}
}
