package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func Categorize(categorizationServer Destination, img string) (resp CategorizationResponse, err error) {
	var body []byte
	body, err = json.Marshal(ImageRequest{Image: img})
	if err != nil {
		return
	}

	var (
		respBody []byte
		code     int
	)
	for i := 0; i < categorizationServer.Retries+1; i++ {
		var req *http.Request
		req, err = http.NewRequest("POST", "http://"+categorizationServer.Address+"/api/categorize", bytes.NewReader(body))

		if err != nil {
			return
		}
		respBody, code, err = RecieveResponseJSON(http.DefaultClient.Do(req))
		if err == nil {
			break
		}
		time.Sleep(categorizationServer.RetryInterval)
	}
	if err != nil {
		return
	}
	if code != http.StatusOK {
		return
	}

	var res CategorizationResponse
	if err = json.Unmarshal(respBody, &res); err != nil {
		return
	}
	return res, nil
}
func GetCamID(categorizationServer Destination, img string) (resp string, err error) {
	var body []byte
	body, err = json.Marshal(ImageRequest{Image: img})
	if err != nil {
		return "", err
	}

	var (
		respBody []byte
		code     int
	)
	for i := 0; i < categorizationServer.Retries+1; i++ {
		var req *http.Request
		req, err = http.NewRequest("POST", "http://"+categorizationServer.Address+"/api/cam-id", bytes.NewReader(body))

		if err != nil {
			return
		}
		respBody, code, err = RecieveResponseJSON(http.DefaultClient.Do(req))
		if err == nil {
			break
		}
		time.Sleep(categorizationServer.RetryInterval)
	}
	if err != nil {
		return
	}
	if code != http.StatusOK {
		return
	}

	var res string
	if err = json.Unmarshal(respBody, &res); err != nil {
		return
	}
	return res, nil
}
