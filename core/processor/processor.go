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

	//for faster test
	imgs = imgs[:10]

	camIDs, err := cv.GetImagesCamIDs(dir, imgs)
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

	addrReqs := make([]database.SetAddressRequest, len(imgs))
	for i := range addrReqs {
		addrReqs[i].Filename = imgs[i]
		addrReqs[i].CamID = camIDs[i]
		addrReqs[i].Address = "addr"
	}
	return database.SetAddress(addrReqs)
}
