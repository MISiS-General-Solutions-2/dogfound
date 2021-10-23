package http

import (
	"io"
	"io/ioutil"
	"net/http"
)

func ClearReponse(resp *http.Response) {
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}

func RecieveResponseJSON(resp *http.Response, err error) ([]byte, int, error) {
	if err != nil {
		ClearReponse(resp)
		return nil, resp.StatusCode, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	resp.Body.Close()
	return b, resp.StatusCode, err
}
