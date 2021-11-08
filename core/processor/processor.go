package processor

import (
	"dogfound/cv"
	"dogfound/database"
	"dogfound/errors"
	"dogfound/http"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

var (
	Processor processor

	counter = 0
	total   = 0
)

type processor struct {
	Config
	cameraInput    chan string
	volunteerInput chan volunteerAddedImage
	addressGuesses chan string
	wg             sync.WaitGroup

	statuses map[string]struct{}
	mu       sync.Mutex

	t1 time.Time
}

func CreateProcessor(cfg *Config) *processor {
	proc := processor{Config: *cfg}
	proc.statuses = make(map[string]struct{})
	proc.cameraInput = make(chan string, 100)
	proc.volunteerInput = make(chan volunteerAddedImage, 1000)
	proc.addressGuesses = make(chan string, 400)
	return &proc
}
func (r *processor) Start() {
	r.start()
	for {
		if err := r.addNewImages(); err != nil {
			fmt.Println(err)
		}
		time.Sleep(r.SampleInterval)
	}
}
func (r *processor) shouldEnqueueImage(img string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.statuses[img]; exists {
		return false
	}
	r.statuses[img] = struct{}{}
	return true
}
func (r *processor) addNewImages() error {
	imgs, err := database.GetNewImages(r.CameraInputDirectory)
	if err != nil {
		return err
	}
	if len(imgs) == 0 {
		return nil
	}

	r.mu.Lock()
	if len(r.statuses) == 0 {
		total = len(imgs)
		fmt.Println("total: ", total)
		r.t1 = time.Now()
	}
	r.mu.Unlock()

	for _, img := range imgs {
		if r.shouldEnqueueImage(img) {
			r.cameraInput <- img
		}
	}
	return nil
}
func (r *processor) start() {
	for i := 0; i < r.NumWorkers; i++ {
		r.wg.Add(1)
		go r.worker()
	}
}
func (r *processor) worker() {
	defer r.wg.Done()
	for {
		select {
		case image := <-r.cameraInput:
			r.process(image)
		case image := <-r.volunteerInput:
			r.processVolunteerAdded(image)
		case image := <-r.addressGuesses:
			r.processAddressGuesses(image)
		}
	}

}
func (r *processor) dropImage(image string) {
	if err := os.Remove(r.CameraInputDirectory + image); err != nil {
		log.Println(err)
	}
	r.mu.Lock()
	delete(r.statuses, image)
	r.mu.Unlock()
}
func (r *processor) process(image string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic %v processing image %v", err, image)
			fmt.Println(string(debug.Stack()))
		}
	}()

	camID, timestamp, err := cv.ParseImage(r.CameraInputDirectory + image)
	if err != nil {
		fmt.Printf("Dropping bad image %v because of error %v\n", image, err)
		r.dropImage(image)
		return
	}
	if camID == "" {
		camID = "PVN_hd_SVAO_3498_4"
	}
	var cr http.CategorizationResponse
	if cr, err = r.GetClassInfo(r.CameraInputDirectory + image); err != nil {
		if !errors.IsDestinationError(err) {
			log.Printf("Dropping bad image %v because of error %v\n", image, err)
			r.dropImage(image)
		} else {
			r.mu.Lock()
			delete(r.statuses, image)
			r.mu.Unlock()
			log.Printf("Destination error occured for image %v. Will try later. Error is %v\n", image, err)
		}
		return
	}
	r.saveProcessedCameraImage(image, camID, timestamp, cr)
	counter += 1
	fmt.Println(counter)
	if counter == total {
		fmt.Println("finished in ", time.Since(r.t1).Seconds())
		counter = 0
		withUnparsedCamIDs, err := database.SelectWithUnparsedCamIDs()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, image := range withUnparsedCamIDs {
			r.addressGuesses <- image
		}
	}
}
func (r *processor) saveProcessedCameraImage(image, camID string, timestamp int64, cr http.CategorizationResponse) {
	record := database.ImagesRecord{
		Filename:  image,
		ClassInfo: cr.ClassInfo,
		CamID:     camID,
		TimeStamp: timestamp,
	}
	defer func() {
		r.mu.Lock()
		delete(r.statuses, image)
		r.mu.Unlock()
	}()
	if err := database.AddAdditionalData(image, cr.Additional); err != nil {
		log.Println(err)
	}
	if err := database.AddImage(r.CameraInputDirectory, record); err != nil {
		log.Println(err)
	}
}

func (r *processor) GetClassInfo(image string) (http.CategorizationResponse, error) {
	if r.Classificator.Address == "" {
		return http.CategorizationResponse{}, nil
	}
	return http.Categorize(r.Classificator, image)
}
func (r *processor) processVolunteerAdded(image volunteerAddedImage) {
	cr, err := r.GetClassInfo(r.VolunteerInputDirectory + image.filename)
	if err != nil {
		if !errors.IsDestinationError(err) {
			log.Printf("Dropping bad image %v because of error %v\n", image, err)
			if err := os.Remove(r.VolunteerInputDirectory + image.filename); err != nil {
				log.Println(err)
			}
		} else {
			r.volunteerInput <- image
			log.Printf("Destination error occured for image %v. Will try later. Error is %v\n", image, err)
		}
		return
	}
	r.saveProcessedVolunteerImage(image, cr)
}
func (r *processor) saveProcessedVolunteerImage(image volunteerAddedImage, cr http.CategorizationResponse) {
	record := database.ImagesRecord{
		Filename:  image.filename,
		ClassInfo: cr.ClassInfo,
		CamID:     "",
		TimeStamp: int64(image.timestamp),
	}
	if err := database.AddAdditionalData(image.filename, cr.Additional); err != nil {
		log.Println(err)
	}
	if err := database.AddImage(r.VolunteerInputDirectory, record); err != nil {
		log.Println(err)
		return
	}
	if err := database.AddVolunteerSourcedAdditonalData(record.Filename, image.lonlat); err != nil {
		log.Println(err)
	}
}
func (r *processor) EnqueueVolunteerImage(image string, timestamp int, lat, lon float64) {
	r.volunteerInput <- volunteerAddedImage{filename: image, timestamp: timestamp, lonlat: [2]float64{lon, lat}}
}
func (r *processor) processAddressGuesses(image string) {
	return
	defaultCamID := ""
	camID, err := http.GetCamID(r.Classificator, image)
	if err != nil {
		fmt.Println(err)
		camID = defaultCamID
	}
	if camID == "" {
		camID = defaultCamID
	}
	if err := database.AddCamID(image, camID); err != nil {
		fmt.Println(err)
	}
}
