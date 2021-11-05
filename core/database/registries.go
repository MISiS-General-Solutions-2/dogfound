package database

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

type registryRecord struct {
	ID       string
	Address  string
	Lat, Lon float64
}

func (r *registryRecord) UnmarshalJSON(b []byte) error {
	type GeoData struct {
		Coordinates [2]float64 `json:"coordinates"`
	}
	type RegistryRecordJSON struct {
		Address string
		ID      string
		Geo     GeoData `json:"geoData"`
	}
	var rr RegistryRecordJSON
	if err := json.Unmarshal(b, &rr); err != nil {
		return err
	}
	r.Address = rr.Address
	r.ID = rr.ID
	r.Lon = rr.Geo.Coordinates[0]
	r.Lat = rr.Geo.Coordinates[1]
	return nil
}
func setRegistryData(recs []registryRecord) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
	INSERT OR IGNORE INTO registries(cam_id,address,lon,lat)
		VALUES(?,?,?,?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, rec := range recs {
		_, err = stmt.Exec(rec.ID, rec.Address, rec.Lon, rec.Lat)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
func updateRegistryAddresses(idAddressPair [][2]string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
	UPDATE registries
	SET
		address=?
	WHERE cam_id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, pair := range idAddressPair {
		_, err = stmt.Exec(pair[1], pair[0])
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
func PopulateRegistries() error {

	if err := loadRegistriesFromJSON(); err != nil {
		return err
	}
	if err := enrichWithAddressesFromXLSX(); err != nil {
		return err
	}
	return nil
}
func loadRegistriesFromJSON() error {
	entries, err := os.ReadDir(registriesPath)
	if err != nil {
		return err
	}
	var records []registryRecord
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".json") {
			b, err := os.ReadFile(registriesPath + entry.Name())
			if err != nil {
				return err
			}
			if err = json.Unmarshal(b, &records); err != nil {
				return err
			}
			if err = setRegistryData(records); err != nil {
				return err
			}
		}
	}
	return nil
}
func enrichWithAddressesFromXLSX() error {
	entries, err := os.ReadDir(registriesPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".xlsx") {
			f, err := excelize.OpenFile(registriesPath + entry.Name())
			if err != nil {
				return err
			}
			rows, err := f.GetRows("0")
			if err != nil {
				return err
			}
			ee := make([][2]string, 0, len(rows)-1)
			for i := 1; i < len(rows); i++ {
				ee = append(ee, [2]string{rows[i][4], rows[i][0]})
			}
			if err = updateRegistryAddresses(ee); err != nil {
				return err
			}
		}
	}
	return nil
}
