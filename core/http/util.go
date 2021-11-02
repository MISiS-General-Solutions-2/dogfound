package http

import (
	"dogfound/errors"
	"io"
	"io/ioutil"
	"net/http"
)

func ClearReponse(resp *http.Response) {
	if resp != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func RecieveResponseJSON(resp *http.Response, err error) ([]byte, int, error) {
	if err != nil {
		return nil, 0, errors.NewDestinationError(err)
	}
	defer ClearReponse(resp)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, errors.NewDestinationError(err)
	}
	return b, resp.StatusCode, nil
}
