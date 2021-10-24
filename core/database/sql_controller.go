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
func GetImageCount() (int, error) {
	q := "SELECT COUNT(color) FROM images WHERE color!=NULL"
	resp, err := db.Query(q)
	if err != nil {
		return 0, err
	}
	defer resp.Close()
	var count sql.NullInt64
	resp.Scan(&count)
	return int(count.Int64), nil
}
func GetNewImages(images []string) (res []string, err error) {

	q := "CREATE TEMP TABLE files(filename TEXT NOT NULL PRIMARY KEY);"
	_, err = db.Exec(q)
	if err != nil {
		return nil, err
	}
	defer func() {
		q := "DROP TABLE files;"
		_, deferErr := db.Exec(q)
		if err != nil {
			err = fmt.Errorf("error during dropping temp table: %v", deferErr)
		}
	}()

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(`INSERT INTO files VALUES(?)`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	for _, img := range images {
		_, err = stmt.Exec(img)
		if err != nil {
			return nil, err
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	q = "SELECT files.filename FROM files WHERE files.filename NOT IN (SELECT images.filename FROM images);"
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	var result []string
	defer rows.Close()
	for rows.Next() {
		var img string
		err = rows.Scan(&img)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, img)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, err
}
func DropRecordsForDeletedImages(images []string) error {
	list := strings.Join(images, `","`)
	stmt := fmt.Sprintf(`DELETE FROM images WHERE filename not in ("%v")`, list)
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
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
	for k := range req {
		switch k {
		case Color, Tail, CamID, TimeStamp:
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
