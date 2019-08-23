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

// func readQuery() {
// 	return baseRequest(method string, node string, params []string)
// }

func RPCRequest(method, node string) string {
	params := [...]string{"0xD7Ff41CA29306122185A07d04293DdB35F24Cf2d", "latest"}
	return baseRequest(method, node, params[:])
}
