package governance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
)

func postAndParse(url string, postData []byte) (map[string]interface{}, error) {
	resp, err := http.Post(string(url), "application/json", bytes.NewReader(postData))
	if err != nil {
		return nil, errors.Wrapf(err, "could not send post request")
	}
	defer resp.Body.Close()
	return parseAndUnmarshal(resp)
}

func parseAndUnmarshal(resp *http.Response) (map[string]interface{}, error) {
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// the 500 is used for errors like `invalid signature`, `not in voting window`, etc.
	if resp.StatusCode != 200 && resp.StatusCode != 500 {
		return nil, errors.Errorf("unexpected response code %s\nbody: %s", resp.Status, bodyData)
	}

	if fastjson.GetString(bodyData, "error") != "" {
		return nil, fmt.Errorf("error: %s, %s", fastjson.GetString(bodyData, "error"), fastjson.GetString(bodyData, "error_description"))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyData, &result); err != nil {
		return nil, errors.Wrapf(err, "could not decode result from %s", string(bodyData))
	} else {
		return result, nil
	}
}
