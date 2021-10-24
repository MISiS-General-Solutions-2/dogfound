package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	nnServiceCfg = Config{
		Address: "neural_network:80",
	}
	ocrServiceCfg = Config{
		Address: "ocr:80",
	}
)

func Categorize(dir string, imgs []string) ([]CategorizationResponse, error) {
	if len(imgs) == 0 {
		return nil, nil
	}
	body, err := json.Marshal(ImageRequest{Dir: dir, Images: imgs})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "http://"+nnServiceCfg.Address+"/api/categorize", bytes.NewReader(body))

	fmt.Println(req.URL)
	if err != nil {
		return nil, err
	}
	var (
		respBody []byte
		code     int
	)
	for i := 0; i < 5; i++ {
		respBody, code, err = RecieveResponseJSON(http.DefaultClient.Do(req))
		if err == nil && code == http.StatusOK {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		return nil, fmt.Errorf("unexpected response code %v with body %s", code, respBody)
	}

	var res []CategorizationResponse
	if err = json.Unmarshal(respBody, &res); err != nil {
		return nil, err
	}
	return res, nil
}
