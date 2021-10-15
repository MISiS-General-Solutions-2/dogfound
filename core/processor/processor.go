package processor

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	phttp "pet-track/http"
)

func GetImageText(cfg *Config, image string) error {
	file, err := os.Open(image)
	if err != nil {
		return err
	}
	resp, err := phttp.Upload(http.DefaultClient, cfg.OCRServerURL, map[string]io.Reader{"file": file})
	resp, err = phttp.RecieveResponse(resp, http.StatusOK, err)
	if err != nil {
		return err
	}
	defer phttp.ClearReponse(resp)

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	log.Printf("%s\n", b)

	return nil
}
