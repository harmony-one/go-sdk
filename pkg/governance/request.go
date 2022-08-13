package governance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/valyala/fastjson"
)

func postAndParse(url string, postData []byte) (map[string]interface{}, error) {
	resp, err := http.Post(string(url), "application/json", bytes.NewReader(postData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseAndUnmarshal(resp)
}

func parseAndUnmarshal(resp *http.Response) (map[string]interface{}, error) {
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if fastjson.GetString(bodyData, "error") != "" {
		return nil, fmt.Errorf("error: %s, %s", fastjson.GetString(bodyData, "error"), fastjson.GetString(bodyData, "error_description"))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyData, &result); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
