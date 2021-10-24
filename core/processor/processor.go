package processor

import (
	"dogfound/cv"
	"dogfound/database"
	"dogfound/http"
	"fmt"
	"math"
	"runtime"
	"sync"
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

		// if err = GetImageClassInfo(dir, imgs); err != nil {
		// 	return err
		// }
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

	numworkers := runtime.GOMAXPROCS(0) - 2
	wg := sync.WaitGroup{}
	camCh := make(chan []string, numworkers)
	timestampsCh := make(chan []int64, numworkers)
	for i := 0; i < numworkers; i++ {
		wg.Add(1)
		var err error
		start, end := getPartOfSlice(len(imgs), i, numworkers)
		go func() {
			defer wg.Done()
			var (
				camIDs     []string
				timestamps []int64
			)
			camIDs, timestamps, err = cv.ParseImages(dir, imgs[start:end])
			if err != nil {
				return
			}
			camCh <- camIDs
			timestampsCh <- timestamps
		}()
	}
	wg.Wait()

	classReqs := make([]database.SetClassesRequest, len(imgs))
	for i := range classReqs {
		classReqs[i].Filename = imgs[i]
	}
	if err := database.SetClasses(classReqs); err != nil {
		return err
	}

	addrReqs := make([]database.CameraInfo, len(imgs))
	i := 0
	for {
		if len(camCh) == 0 {
			break
		}
		camIDs := <-camCh
		timestamps := <-timestampsCh
		addrReqs[i].Filename = imgs[i]
		addrReqs[i].CamID = camIDs[i]
		addrReqs[i].TimeStamp = timestamps[i]
		i += 1
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
