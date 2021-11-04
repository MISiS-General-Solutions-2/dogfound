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
	counter = 0
	total   = 0
)

type processor struct {
	Config
	inputChannel chan string
	wg           sync.WaitGroup

	statuses map[string]struct{}
	mu       sync.Mutex

	t1 time.Time
}

func StartProcessor(cfg *Config) error {
	proc := processor{Config: *cfg}
	proc.statuses = make(map[string]struct{})
	proc.inputChannel = make(chan string, 100)
	proc.start()
	for {
		if err := proc.addNewImages(); err != nil {
			return err
		}
		time.Sleep(proc.SampleInterval)
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
	imgs, err := database.GetNewImages(r.ImageSourceDirectory)
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
			r.inputChannel <- img
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
	for image := range r.inputChannel {
		r.process(image)
	}
}
func (r *processor) dropImage(image string) {
	if err := os.Remove(r.ImageSourceDirectory + image); err != nil {
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

	camID, timestamp, err := cv.ParseImage(r.ImageSourceDirectory + image)
	if err != nil {
		fmt.Printf("Dropping bad image %v because of error %v\n", image, err)
		r.dropImage(image)
		return
	}
	var classInfo database.ClassInfo
	if classInfo, err = r.GetClassInfo(image); err != nil {
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
	r.saveProcessedImage(image, camID, timestamp, classInfo)
	counter += 1
	fmt.Println(counter)
	if counter == total {
		fmt.Println("finished in ", time.Since(r.t1).Seconds())
	}
}
func (r *processor) saveProcessedImage(image, camID string, timestamp int64, classInfo database.ClassInfo) {
	record := database.ImagesRecord{
		Filename:  image,
		ClassInfo: classInfo,
		CamID:     camID,
		TimeStamp: timestamp,
	}
	defer func() {
		r.mu.Lock()
		delete(r.statuses, image)
		r.mu.Unlock()
	}()
	if err := database.AddImage(r.ImageSourceDirectory, record); err != nil {
		log.Println(err)
	}
}

func (r *processor) GetClassInfo(img string) (database.ClassInfo, error) {
	if r.Classificator.Address == "" {
		return database.ClassInfo{}, nil
	}
	return http.Categorize(r.Classificator, r.ImageSourceDirectory+img)
}
