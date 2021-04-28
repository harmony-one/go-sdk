package governance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/valyala/fastjson"
)

func getAndParse(url governanceApi, data interface{}) error {
	resp, err := http.Get(string(url))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return parseAndUnmarshal(resp, data)
}

func postAndParse(url governanceApi, postData []byte, data interface{}) error {
	resp, err := http.Post(string(url), "application/json", bytes.NewReader(postData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return parseAndUnmarshal(resp, data)
}

func parseAndUnmarshal(resp *http.Response, data interface{}) error {
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if fastjson.GetString(bodyData, "error") != "" {
		return fmt.Errorf("error: %s, %s", fastjson.GetString(bodyData, "error"), fastjson.GetString(bodyData, "error_description"))
	}

	return json.Unmarshal(bodyData, data)
}
