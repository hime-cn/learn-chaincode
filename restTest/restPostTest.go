package main

import (
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

// rpcResponse defines the JSON RPC 2.0 response payload for the /chaincode endpoint.
type rpcResponse struct {
	Jsonrpc string     `json:"jsonrpc,omitempty"`
	Result  *rpcResult `json:"result,omitempty"`
	Error   *rpcError  `json:"error,omitempty"`
	ID      *rpcID     `json:"id"`
}

// rpcResult defines the structure for an rpc sucess/error result message.
type rpcResult struct {
	Status  string    `json:"status,omitempty"`
	Message string    `json:"message,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
}

// rpcError defines the structure for an rpc error.
type rpcError struct {
	// A Number that indicates the error type that occurred. This MUST be an integer.
	Code int64 `json:"code,omitempty"`
	// A String providing a short description of the error. The message SHOULD be
	// limited to a concise single sentence.
	Message string `json:"message,omitempty"`
	// A Primitive or Structured value that contains additional information about
	// the error. This may be omitted. The value of this member is defined by the
	// Server (e.g. detailed error information, nested errors etc.).
	Data string `json:"data,omitempty"`
}
type rpcID struct {
	StringValue *string
	IntValue    *int64
}

var loc_server1 = "http://192.168.1.242:7050"

func main() {
	var chaincodeid = "964c8793cdc18fd4389626c991b40fc50bb009d63632bdbb57813e35476c97f4a975841e4c3f94aa3e0f48b10ec8697db8be84b7a3c97ff685f1adc932198137"
	for i := 0; i <= 100; i++ {
		httpResponse, body := postInvoke(strconv.Itoa(i), chaincodeid, "567")
		if httpResponse.StatusCode != http.StatusOK {
			fmt.Errorf("Expected an HTTP status code %#v but got %#v", http.StatusOK, httpResponse.StatusCode)
		}
		res := parseRPCResponse(body)
		if res.Error == nil {
			fmt.Printf("res.Result.Message:%s", res.Result.Message)
		} else {
			fmt.Printf("res.Error.Message:%s", res.Error.Message)
		}
	}
}

func postInvoke(id string, chaincodeid string, acctno string) (*http.Response, []byte) {
	var posturl = loc_server1 + "/chaincode";
	requestBody := `{
		"jsonrpc": "2.0",
		"ID": ` +
		id +
		`,
				"method": "invoke",
				"params": {
					"type": 1,
					"chaincodeID": {
						"name": "` +
		chaincodeid +
		`"
					},
					"ctorMsg": {
						 "function": "createAccount",
						"args": ["` +
		acctno +
		`","acct` +
		id +
		`","100"]
			},
			"secureContext": "admin"
		}
		}`
	response, err := http.Post(posturl, "application/json", bytes.NewReader([]byte(requestBody)))
	if err != nil {
		fmt.Errorf("Error attempt to POST %s: %v", posturl, err)
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		fmt.Errorf("Error reading HTTP resposne body: %v", err)
	}
	return response, body
}

func parseRPCResponse(body []byte) rpcResponse {
	var res rpcResponse
	err := json.Unmarshal(body, &res)
	if err != nil {
		fmt.Errorf("Invalid JSON RPC response: %v", err)
	}
	return res
}
