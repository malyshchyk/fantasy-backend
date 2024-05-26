package handlers

import (
	"io"
	"net/http"
	"net/url"
	"sync"
)

func GetAsync(baseURL string, params url.Values, ch chan<- []byte, errch chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	u, err := url.Parse(baseURL)
	if err != nil {
		errch <- err
		return
	}
	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		errch <- err
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errch <- err
		return
	}
	ch <- body
}

func Get(baseURL string, params url.Values) ([]byte, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return body, err
}
