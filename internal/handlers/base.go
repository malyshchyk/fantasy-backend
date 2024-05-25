package handlers

import (
	"io"
	"net/http"
	"net/url"
)

func Get(baseURL string, params url.Values, w http.ResponseWriter) []byte {
	u, err := url.Parse(baseURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	return body
}
