package processor

import (
	"fmt"
	"pet-track/cv"
	"pet-track/database"
	"pet-track/http"
)

func ProcessAllImages(cfg *Config) error {
	dir, imgs, err := database.GetImages()
	if err != nil {
		return err
	}

	if err = GetImageClassInfo(dir, imgs); err != nil {
		return err
	}
	return GetOCRInfo(dir, imgs)
}
func GetOCRInfo(dir string, imgs []string) error {
	camIDs, err := cv.GetImagesCamIDs(dir, imgs)
	if err != nil {
		return err
	}
	timestamps, err := GetTimestampsMock(dir, imgs)
	if err != nil {
		return err
	}

	classReqs := make([]database.SetClassesRequest, len(imgs))
	for i := range classReqs {
		classReqs[i].Filename = imgs[i]
	}
	if err := database.SetClasses(classReqs); err != nil {
		return err
	}

	addrReqs := make([]database.CameraInfo, len(imgs))
	for i := range addrReqs {
		addrReqs[i].Filename = imgs[i]
		addrReqs[i].CamID = camIDs[i]
		addrReqs[i].TimeStamp = timestamps[i]
	}
	return database.SetCameraInfo(addrReqs)
}
func GetImageClassInfo(dir string, imgs []string) error {
	res, err := http.Categorize(dir, imgs)
	if err != nil {
		return err
	}
	for i := range res {
		fmt.Println(res[i].Vis.Probabilities)
	}
	return nil
}
func GetTimestampsMock(dir string, imgs []string) ([]int64, error) {
	return make([]int64, len(imgs)), nil
}
