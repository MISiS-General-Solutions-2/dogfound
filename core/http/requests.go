package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Categorize(cfg Config, dir string, imgs []string) ([]CategorizationResponse, error) {
	body, err := json.Marshal(ImageRequest{Dir: dir, Images: imgs})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", cfg.Address+"/api/recognize", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	respBody, code, err := RecieveResponseJSON(http.DefaultClient.Do(req))
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
