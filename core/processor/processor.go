package processor

import (
	"dogfound/cv"
	"dogfound/database"
	"dogfound/http"
	"log"
	"math"
	"time"
)

// blocks
func ProcessNewImages() (err error) {
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
			log.Println(err)
		}
		if err = GetOCRTextInfo(dir, imgs); err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}
}
func getPartOfSlice(l, i, numpatrs int) (int, int) {
	start := (l / numpatrs) * i
	end := int(math.Min(float64(l), float64(l/numpatrs*(i+1))))
	return start, end
}
func GetOCRTextInfo(dir string, imgs []string) error {
	if len(imgs) == 0 {
		return nil
	}
	imgs = imgs[:10]

	camIDs, timestamps, err := cv.ParseImages(dir, imgs)
	if err != nil {
		return err
	}

	addrReqs := make([]database.CameraInfo, len(imgs))
	for i := range imgs {
		addrReqs[i].Filename = imgs[i]
		addrReqs[i].CamID = camIDs[i]
		addrReqs[i].TimeStamp = timestamps[i]
	}

	return database.SetCameraInfo(addrReqs)
}
func GetImageClassInfo(dir string, imgs []string) error {
	i := 10
	for {
		if i >= len(imgs) {
			break
		}
		res, err := http.Categorize(dir, imgs[:i])
		if err != nil {
			return err
		}
		if err = database.SetClasses(res); err != nil {
			return err
		}

		imgs = imgs[i:]
		i += 10
	}
	res, err := http.Categorize(dir, imgs)
	if err != nil {
		return err
	}
	if err = database.SetClasses(res); err != nil {
		return err
	}
	return nil
}
