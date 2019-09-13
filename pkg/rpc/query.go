package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/harmony-one/go-sdk/pkg/common"
)

var (
	queryID = 0
)

func baseRequest(method RPCMethod, node string, params interface{}) []byte {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": common.JSONRPCVersion,
		"id":      strconv.Itoa(queryID),
		"method":  method,
		"params":  params,
	})
	resp, err := http.Post(node, "application/json", bytes.NewBuffer(requestBody))
	if common.DebugRPC {
		fmt.Printf("URL: %s, Request Body: %s\n\n", node, string(requestBody))
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer resp.Body.Close()
	// TODO Need to read the body, fail with the error values that come from ErrorCode list
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	queryID++
	if common.DebugRPC {
		fmt.Printf("URL: %s, Response Body: %s\n\n", node, string(body))
	}
	return body
}

// TODO add the error code usage here, change return signature, make CLI be consumer that checks error
func Request(method RPCMethod, node string, params interface{}) (map[string]interface{}, error) {
	rpcJson := make(map[string]interface{})
	rawReply := baseRequest(method, node, params)
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

func RawRequest(method RPCMethod, node string, params interface{}) []byte {
	return baseRequest(method, node, params)
}
