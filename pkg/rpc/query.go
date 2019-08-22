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
	queryID = 0
)

func baseRequest(method, node string, params []string) string {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": JSON_RPC_VERSION,
		"id":      strconv.Itoa(queryID),
		"method":  method,
		"params":  params,
	})

	resp, err := http.Post(node, "application/json", bytes.NewBuffer(requestBody))
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

func RPCRequest(method, node string) string {
	params := [...]string{"0xD7Ff41CA29306122185A07d04293DdB35F24Cf2d", "latest"}
	return baseRequest(node, method, params[:])
}
