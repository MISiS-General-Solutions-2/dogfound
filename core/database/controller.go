package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
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

	imgPath = "../data/img/"
)

type Record struct {
	Img   string
	Color int
	Tail  int
}

func GetFilePath(name string) string {
	return imgPath + name
}
func Connect() func() {
	var err error
	db, err = sql.Open("sqlite3", "../data/features.db")
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
		result[i] = f.Name()
	}
	return result, nil
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
