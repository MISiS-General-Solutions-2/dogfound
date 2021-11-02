package http

import (
	"bytes"
	"dogfound/database"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	nnServiceCfg = Config{
		Address: "neural_network:80",
	}
)

func Categorize(dir string, imgs []string) ([]database.SetClassesRequest, error) {
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
	respBody, code, err = RecieveResponseJSON(http.DefaultClient.Do(req))
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		return nil, fmt.Errorf("unexpected response code %v with body %s", code, respBody)
	}

	var res []database.SetClassesRequest
	if err = json.Unmarshal(respBody, &res); err != nil {
		return nil, err
	}
	for i := range imgs {
		res[i].Filename = imgs[i]
	}
	return res, nil
}
