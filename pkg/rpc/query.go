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

func baseRequest(method RPCMethod, node string, params []string) string {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": JSON_RPC_VERSION,
		"id":      strconv.Itoa(queryID),
		"method":  method,
		"params":  params,
	})

	resp, err := http.Post(node, "application/json", bytes.NewBuffer(requestBody))
	if debugEnabled {
		fmt.Printf("URL: %s, Body: %s\n", node, string(requestBody))
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	queryID++
	return string(body)
}

func RPCRequest(method, node string, params []string) string {
	return baseRequest(method, node, params)
}
