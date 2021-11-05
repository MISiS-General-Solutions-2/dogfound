package database

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

	LatLon = "latlon"
)

type ImagesRecord struct {
	Filename string `json:"filename"`

	ClassInfo `json:",inline"`

	CamID     string `json:"cam_id"`
	TimeStamp int64  `json:"timestamp"`
}

type ClassInfo struct {
	IsAnimal    int `json:"is_animal_there"`
	IsDog       int `json:"is_it_a_dog"`
	IsWithOwner int `json:"is_the_owner_there"`
	Color       int `json:"color"`
	Tail        int `json:"tail"`
}
type CameraInfo struct {
	Filename string `json:"filename"`

	CamID     string `json:"cam_id"`
	TimeStamp int64  `json:"timestamp"`
}
type SearchResponse struct {
	Filename  string     `json:"filename"`
	Address   string     `json:"address"`
	CamID     string     `json:"cam_id"`
	TimeStamp int64      `json:"timestamp"`
	LonLat    [2]float64 `json:"lonlat"`

	AdditionalData Additional `json:"additional"`
}
type Additional struct {
	Crop [4]int `json:"crop"`
}
