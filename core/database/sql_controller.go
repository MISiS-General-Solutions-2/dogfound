package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

const (
	Filename    = "filename"
	IsAnimal    = "is_animal_there"
	IsDog       = "is_it_a_dog"
	IsWithOwner = "is_the_owner_there"
	Color       = "color"
	Tail        = "tail"
	Address     = "address"
	CamID       = "cam_id"
	TimeStamp   = "timestamp"

	DataPath = "/opt/pet-track/data/"
)

func Connect() func() {
	var err error
	db, err = sql.Open("sqlite3",
		DataPath+"images.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS images
	(
		filename TEXT NOT NULL PRIMARY KEY,
		is_animal_there INTEGER,
		is_it_a_dog INTEGER,
		is_the_owner_there INTEGER,
		color INTEGER,
		tail INTEGER,
		address TEXT,
		cam_id TEXT,
		timestamp INTEGER
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

type SetClassesRequest struct {
	Filename string

	IsAnimal    int
	IsDog       int
	IsWithOwner int
	Color       int
	Tail        int
}
type SetAddressRequest struct {
	Filename string

	Address   string
	CamID     string
	TimeStamp int
}

func SetAddress(reqs []SetAddressRequest) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
	INSERT INTO images(filename,address,cam_id,timestamp)
	VALUES(?,?,?,?)
	ON CONFLICT(filename)
	DO UPDATE
	SET
	address=?,
	cam_id=?,
	timestamp=?
	WHERE filename = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, rec := range reqs {
		if !imageExists(rec.Filename) {
			return fmt.Errorf("image %v does not exist", rec.Filename)
		}
		_, err = stmt.Exec(rec.Filename, rec.Address, rec.CamID, rec.TimeStamp, rec.Address, rec.CamID, rec.TimeStamp, rec.Filename)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func SetClasses(reqs []SetClassesRequest) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
	INSERT INTO images(filename,is_animal_there,is_the_owner_there,color,tail)
	VALUES(?,?,?,?,?)
	ON CONFLICT(filename)
	DO UPDATE
	SET
	is_animal_there=?,
	is_the_owner_there=?,
	color=?,
	tail=?
	WHERE filename = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, rec := range reqs {
		if !imageExists(rec.Filename) {
			return fmt.Errorf("image %v does not exist", rec.Filename)
		}
		_, err = stmt.Exec(rec.Filename, rec.IsAnimal, rec.IsWithOwner, rec.Color, rec.Tail, rec.IsAnimal, rec.IsWithOwner, rec.Color, rec.Tail, rec.Filename)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
func ValidateRequest(req map[string]interface{}) error {
	for k := range req {
		switch k {
		case IsAnimal, IsWithOwner, Color, Tail, Address, CamID, TimeStamp:
		default:
			return fmt.Errorf("unexpected field %v", k)
		}
	}
	return nil
}

type SearchResponse struct {
	Filename  string
	Address   string
	CamID     string
	TimeStamp int
}

func GetImagesByClasses(req map[string]interface{}) ([]SearchResponse, error) {
	b := strings.Builder{}
	b.WriteString("SELECT filename,address,cam_id,timestamp FROM images WHERE ")
	for k, v := range req {
		b.WriteString(k)
		b.WriteRune('=')
		switch k {
		case Filename, Address, CamID:
			b.WriteRune('"')
			b.WriteString(v.(string))
			b.WriteRune('"')
		case IsAnimal, IsDog, IsWithOwner, Color, Tail:
			b.WriteString(strconv.Itoa(int(v.(float64))))
		default:
			return nil, fmt.Errorf("unexpected field %v", k)
		}
		b.WriteString(" AND ")
	}
	b.WriteString("address IS NOT NULL AND cam_id IS NOT NULL AND timestamp IS NOT NULL")
	sqlStmt := b.String()
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []SearchResponse
	for rows.Next() {
		var sr SearchResponse
		err = rows.Scan(&sr.Filename, &sr.Address, &sr.CamID, &sr.TimeStamp)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, sr)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}
