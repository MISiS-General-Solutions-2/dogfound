package database

import (
	"database/sql"
	"errors"
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
	DataPath       = "/opt/pet-track/data/"
	registriesPath = DataPath + "registries/"
)

func Connect() func() {
	var err error
	db, err = sql.Open("sqlite3",
		DataPath+"pet-track.db")
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
	for k, v := range req {
		switch k {
		case Color, Tail, CamID, TimeStamp:
		case LatLon:
			if arr, ok := v.([]interface{}); ok && len(arr) == 2 {
				for i := range arr {
					if _, ok := arr[i].(float64); !ok {
						return errors.New("latlon should be array of two floats")
					}
				}
			} else {
				return errors.New("latlon should be array of two floats")
			}
		default:
			return fmt.Errorf("unexpected field %v", k)
		}
	}
	return nil
}

func GetImagesByClasses(req map[string]interface{}) ([]SearchResponse, error) {
	b := strings.Builder{}
	b.WriteString(`SELECT filename,registries.address,images.cam_id,lat,lon,timestamp FROM images LEFT OUTER JOIN registries
		ON images.cam_id = registries.cam_id WHERE `)
	first := true
	for k, v := range req {
		switch k {
		case Filename, Address, CamID:
			if !first {
				b.WriteString(" AND ")
			}
			b.WriteString(k)
			b.WriteRune('=')
			b.WriteRune('"')
			b.WriteString(v.(string))
			b.WriteRune('"')
		case IsAnimal, IsDog, IsWithOwner, Color, Tail:
			if !first {
				b.WriteString(" AND ")
			}
			b.WriteString(k)
			b.WriteRune('=')
			b.WriteString(strconv.Itoa(int(v.(float64))))
		case TimeStamp:
			if !first {
				b.WriteString(" AND ")
			}
			b.WriteString("(timestamp=0 OR timestamp>=")
			b.WriteString(strconv.Itoa(int(v.(float64))))
			b.WriteRune(')')
		case LatLon:
			//used later
		default:
			return nil, fmt.Errorf("unexpected field %v", k)
		}
		first = false
	}
	sqlStmt := b.String()
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []SearchResponse
	for rows.Next() {
		type SearchResponseSQL struct {
			Filename  string
			Address   sql.NullString
			CamID     sql.NullString
			TimeStamp sql.NullInt64
			Lat, Lon  sql.NullFloat64
		}
		var sr SearchResponseSQL
		err = rows.Scan(&sr.Filename, &sr.Address, &sr.CamID, &sr.Lat, &sr.Lon, &sr.TimeStamp)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, SearchResponse{sr.Filename, sr.Address.String, sr.CamID.String, sr.TimeStamp.Int64, [2]float64{sr.Lon.Float64, sr.Lat.Float64}, Visualization{}})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}
