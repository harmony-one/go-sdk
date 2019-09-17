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

func baseRequest(method RPCMethod, node string, params interface{}) ([]byte, error) {
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
	fasthttp.ReleaseRequest(req)
	body := res.Body()
	result := make([]byte, len(body))
	copy(result, body)
	fasthttp.ReleaseResponse(res) // Only when you are done with body!
	if common.DebugRPC {
		fmt.Printf("URL: %s, Request Body: %s\n\n", node, string(requestBody))
	}
	if common.DebugRPC {
		fmt.Printf("URL: %s, Response Body: %s\n\n", node, string(body))
	}
	queryID++
	return result, nil
}

// TODO add the error code usage here, change return signature, make CLI be consumer that checks error
func Request(method RPCMethod, node string, params interface{}) (map[string]interface{}, error) {
	rpcJson := make(map[string]interface{})
	rawReply, err := baseRequest(method, node, params)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(rawReply, &rpcJson)
	if oops := rpcJson["error"]; oops != nil {
		errNo := oops.(map[string]interface{})["code"].(float64)
		errMessage := ""
		if oops.(map[string]interface{})["message"] != nil {
			errMessage = oops.(map[string]interface{})["message"].(string)
		}
		return nil, ErrorNumberToError(errMessage, errNo)
	}
	return rpcJson, nil
}

func RawRequest(method RPCMethod, node string, params interface{}) ([]byte, error) {
	return baseRequest(method, node, params)
}
