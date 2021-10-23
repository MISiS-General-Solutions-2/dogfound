package processor

import (
	"pet-track/cv"
	"pet-track/database"
)

func ProcessAllImages(cfg *Config) error {
	dir, imgs, err := database.GetImages()
	if err != nil {
		return err
	}

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
func GetTimestampsMock(dir string, imgs []string) ([]int64, error) {
	return make([]int64, len(imgs)), nil
}
