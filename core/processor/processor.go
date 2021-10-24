package processor

import (
	"dogfound/cv"
	"dogfound/database"
	"dogfound/http"
	"fmt"
	"time"
)

// blocks
func ProcessNewImages() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from %v", r)
		}
	}()
	for {
		dir, imgs, err := database.GetImages()
		if err != nil {
			return err
		}
		count, err := database.GetImageCount()
		if err != nil {
			return err
		}
		if count != len(dir) {
			if err = database.DropRecordsForDeletedImages(imgs); err != nil {
				return err
			}

			imgs, err = database.GetNewImages(imgs)
			if err != nil {
				return err
			}
		}

		if err = GetImageClassInfo(dir, imgs); err != nil {
			return err
		}
		if err = GetOCRInfo(dir, imgs); err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}
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
