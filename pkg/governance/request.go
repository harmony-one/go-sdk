package governance

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func getAndParse(url governanceApi, data interface{}) error {
	resp, err := http.Get(string(url))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(data)
}

func postAndParse(url governanceApi, postData []byte, data interface{}) error {
	resp, err := http.Post(string(url), "application/json", bytes.NewReader(postData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(data)
}
