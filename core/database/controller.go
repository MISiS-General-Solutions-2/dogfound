package database

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

const (
	FeatureColor = "color"
	FeatureTail  = "tail"

	DataPath  = "/opt/pet-track/data/"
	imagePath = DataPath + "img/"
)

type Record struct {
	Img   string
	Color int
	Tail  int
}

func GetFilePath(name string) string {
	return imagePath + name
}
func Connect() func() {
	var err error
	db, err = sql.Open("sqlite3",
		DataPath+"features.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS features
	(
		img TEXT NOT NULL PRIMARY KEY,
		color INTEGER,
		tail INTEGER
	)
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		db.Close()
		panic(err)
	}

	return func() {
		db.Close()
	}
}
func GetFilesInDirectory(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(files))
	for i, f := range files {
		result[i] = dir + f.Name()
	}
	return result, nil
}
func GetImages() []string {
	imgs, err := GetFilesInDirectory(imagePath)
	if err != nil {
		panic(err)
	}
	return imgs
}
func AddImages(imgs []string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT OR IGNORE INTO features(img) VALUES(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, img := range imgs {
		_, err = stmt.Exec(img)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
func SetFeatures(records []Record) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
	UPDATE features
	SET color = ?,
		tail = ?
	WHERE img = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, rec := range records {
		_, err = stmt.Exec(rec.Color, rec.Tail, rec.Img)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
func GetImagesByFeatures(features map[string]interface{}) ([]string, error) {
	b := strings.Builder{}
	b.WriteString("SELECT img FROM features WHERE ")
	i := 0
	for k, v := range features {
		b.WriteString(k)
		b.WriteRune('=')
		switch t := v.(type) {
		case string:
			b.WriteString(v.(string))
		case int:
			b.WriteString(strconv.Itoa(v.(int)))
		default:
			fmt.Printf("unexpected type%v\n", t)
			b.WriteString(fmt.Sprint(v))
		}
		if i < len(features)-1 {
			b.WriteRune(',')
		}
	}
	sqlStmt := b.String()
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []string
	for rows.Next() {
		var img string
		err = rows.Scan(&img)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, img)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}
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
