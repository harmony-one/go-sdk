package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"

	"github.com/harmony-one/go-sdk/pkg/common"
)

var (
	queryID = 0
	post    = []byte("POST")
)

func baseRequest(method string, node string, params interface{}) ([]byte, error) {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": common.JSONRPCVersion,
		"id":      strconv.Itoa(queryID),
		"method":  method,
		"params":  params,
	})
	const contentType = "application/json"
	req := fasthttp.AcquireRequest()
	req.SetBody(requestBody)
	req.Header.SetMethodBytes(post)
	req.Header.SetContentType(contentType)
	req.SetRequestURIBytes([]byte(node))
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		return nil, err
	}
	c := res.StatusCode()
	if c != 200 {
		return nil, fmt.Errorf("http status code not 200, received: %d", c)
	}
	fasthttp.ReleaseRequest(req)
	body := res.Body()
	result := make([]byte, len(body))
	copy(result, body)
	fasthttp.ReleaseResponse(res)
	if common.DebugRPC {
		reqB := common.JSONPrettyFormat(string(requestBody))
		respB := common.JSONPrettyFormat(string(body))
		fmt.Printf("URL: %s, Request Body: %s\n\n", node, reqB)
		fmt.Printf("URL: %s, Response Body: %s\n\n", node, respB)
	}
	queryID++
	return result, nil
}

// TODO Check if Method known, return error when not known, good intern task

// Request processes
func Request(method string, node string, params interface{}) (Reply, error) {
	rpcJSON := make(map[string]interface{})
	rawReply, err := baseRequest(method, node, params)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(rawReply, &rpcJSON)
	if oops := rpcJSON["error"]; oops != nil {
		errNo := oops.(map[string]interface{})["code"].(float64)
		errMessage := ""
		if oops.(map[string]interface{})["message"] != nil {
			errMessage = oops.(map[string]interface{})["message"].(string)
		}
		return nil, ErrorCodeToError(errMessage, errNo)
	}
	return rpcJSON, nil
}

// RawRequest is to sidestep the lifting done by Request
func RawRequest(method string, node string, params interface{}) ([]byte, error) {
	return baseRequest(method, node, params)
}
