package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const (
	JSON_RPC_VERSION = "2.0"
)

var (
	queryID      = 0
	debugEnabled = false
)

func init() {
	if _, enabled := os.LookupEnv("HMY_RPC_DEBUG"); enabled != false {
		debugEnabled = true
	}
}

func baseRequest(method RPCMethod, node string, params interface{}) []byte {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": JSON_RPC_VERSION,
		"id":      strconv.Itoa(queryID),
		"method":  method,
		"params":  params,
	})
	resp, err := http.Post(node, "application/json", bytes.NewBuffer(requestBody))
	if debugEnabled {
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
	if debugEnabled {
		fmt.Printf("URL: %s, Response Body: %s\n\n", node, string(body))
	}
	return body
}

// TODO add the error code usage here, change return signature, make CLI be consumer that checks error
func RPCRequest(method RPCMethod, node string, params interface{}) map[string]interface{} {
	rpcJson := make(map[string]interface{})
	rawReply := baseRequest(method, node, params)
	json.Unmarshal(rawReply, &rpcJson)
	return rpcJson
}
