package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

const (
	DataPath       = "/opt/dogfound/data/"
	registriesPath = DataPath + "registries/"
)

func Connect() func() {
	var err error
	db, err = sql.Open("sqlite3",
		DataPath+"dogfound.db")
	if err != nil {
		log.Fatal(err)
	}
	return func() {
		db.Close()
	}
}
func SetCameraInfo(reqs []CameraInfo) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
	INSERT INTO images(filename,cam_id,timestamp)
		VALUES(?,?,?)
	ON CONFLICT(filename)
	DO UPDATE
	SET
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
		_, err = stmt.Exec(rec.Filename, rec.CamID, rec.TimeStamp, rec.CamID, rec.TimeStamp, rec.Filename)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
func AddAdditionalData(image string, data Additional) error {
	stmt, err := db.Prepare(`INSERT INTO additional(filename,crop_x0,crop_y0,crop_x1,crop_y1)
			VALUES(?,?,?,?,?)
		ON CONFLICT(filename)
		DO UPDATE
		SET
			crop_x0=?,
			crop_y0=?,
			crop_x1=?,
			crop_y1=?
		WHERE filename = ?
		`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(image, data.Crop[0], data.Crop[1], data.Crop[2], data.Crop[3], data.Crop[0], data.Crop[1], data.Crop[2], data.Crop[3], image)
	if err != nil {
		return err
	}
	return nil
}
func AddImage(imageSourceDirectory string, record ImagesRecord) error {
	stmt, err := db.Prepare(`INSERT INTO images(filename,is_animal_there,is_it_a_dog,is_the_owner_there,color,tail,cam_id,timestamp,breed)
			VALUES(?,?,?,?,?,?,?,?,?)
		ON CONFLICT(filename)
		DO UPDATE
		SET
			is_animal_there=?,
			is_it_a_dog=?,
			is_the_owner_there=?,
			color=?,
			tail=?,
			cam_id=?,
			timestamp=?,
			breed=?
		WHERE filename = ?
		`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(record.Filename, record.IsAnimal, record.IsDog, record.IsWithOwner, record.Color, record.Tail, record.CamID, record.TimeStamp, record.Breed, record.IsAnimal, record.IsDog, record.IsWithOwner, record.Color, record.Tail, record.CamID, record.TimeStamp, record.Breed, record.Filename)
	if err != nil {
		return err
	}

	if err = os.Rename(imageSourceDirectory+record.Filename, imagePath+record.Filename); err != nil {
		return err
	}

	return nil
}

func ValidateRequest(req map[string]interface{}) error {
	for k := range req {
		switch k {
		case Color, Tail, CamID, T1, T0:
		default:
			return fmt.Errorf("unexpected field %v", k)
		}
	}
	return nil
}

func addConditions(b *strings.Builder, req map[string]interface{}) error {
	b.WriteString(" WHERE ")
	first := true
	for k, v := range req {
		if !first {
			b.WriteString(" AND ")
		}
		switch k {
		case Filename, Address, CamID:
			b.WriteString(k)
			b.WriteRune('=')
			b.WriteRune('"')
			b.WriteString(v.(string))
			b.WriteRune('"')
		case IsAnimal, IsDog, IsWithOwner, Color, Tail:
			b.WriteString(k)
			b.WriteRune('=')
			b.WriteString(strconv.Itoa(int(v.(float64))))
		case T0:
			b.WriteString("(timestamp=0 OR timestamp>=")
			b.WriteString(strconv.Itoa(int(v.(float64))))
			b.WriteRune(')')
		case T1:
			b.WriteString("(timestamp=0 OR timestamp<=")
			b.WriteString(strconv.Itoa(int(v.(float64))))
			b.WriteRune(')')
		default:
			return fmt.Errorf("unexpected field %v", k)
		}
		first = false
	}
	return nil
}
func GetAdditionalInfo(image string) (*Additional, error) {
	stmt, err := db.Prepare(`SELECT crop_x0,crop_y0,crop_x1,crop_y1
		FROM additional
		WHERE filename=?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(image)
	type AdditionalSQL struct {
		x0, y0, x1, y1 sql.NullInt64
	}
	var a AdditionalSQL

	if err := row.Scan(&a.x0, &a.y0, &a.x1, &a.y1); err != nil {
		if err == sql.ErrNoRows {
			return &Additional{}, nil
		}
		return nil, err
	}
	return &Additional{Crop: [4]int{int(a.x0.Int64), int(a.y0.Int64), int(a.x1.Int64), int(a.y1.Int64)}}, nil
}
func GetImagesByClasses(req map[string]interface{}) ([]SearchResponse, error) {
	b := strings.Builder{}
	b.WriteString(`SELECT images.filename,registries.address,images.cam_id,lat,lon,timestamp,crop_x0,crop_y0,crop_x1,crop_y1,breed FROM images 
		LEFT OUTER JOIN registries 
		ON images.cam_id = registries.cam_id
		LEFT OUTER JOIN additional on images.filename=additional.filename`)
	if len(req) != 0 {
		if err := addConditions(&b, req); err != nil {
			return nil, err
		}
	}
	sqlStmt := b.String()
	fmt.Println(sqlStmt)
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []SearchResponse
	for rows.Next() {
		type SearchResponseSQL struct {
			Filename       string
			Address        sql.NullString
			CamID          sql.NullString
			TimeStamp      sql.NullInt64
			Lat, Lon       sql.NullFloat64
			x0, y0, x1, y1 sql.NullInt64
			Breed          sql.NullString
		}
		var sr SearchResponseSQL
		err = rows.Scan(&sr.Filename, &sr.Address, &sr.CamID, &sr.Lat, &sr.Lon, &sr.TimeStamp, &sr.x0, &sr.y0, &sr.x1, &sr.y1, &sr.Breed)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, SearchResponse{sr.Filename, sr.Address.String, sr.CamID.String, sr.TimeStamp.Int64, [2]float64{sr.Lon.Float64, sr.Lat.Float64}, sr.Breed.String, Additional{Crop: [4]int{int(sr.x0.Int64), int(sr.y0.Int64), int(sr.x1.Int64), int(sr.y1.Int64)}}})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}
