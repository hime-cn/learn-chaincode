package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/hyperledger/fabric/protos"
)

var loc_server = "http://192.168.1.242:7050"

func main() {
	trans, err := getTransactions("c70ff097-5a59-4dc3-b6cc-ac08e8ed6729");
	if err != nil {
		fmt.Printf("err:%s", err)
	}
	fmt.Printf("trans:%s", trans.String())

	getBlocks(1)
}

func getBlocks(i int) {
	var getblockurl = loc_server + "/chain/blocks/" + strconv.Itoa(i);
	var b []byte
	if response, err := http.Get(getblockurl); err != nil {
		return
	} else {
		if b, err = ioutil.ReadAll(response.Body); err != nil {
			return
		}
	}
	var block0 protos.Block
	err := json.Unmarshal(b, &block0)
	if err == nil {
		fmt.Println(block0.String())
		trans := block0.GetTransactions()
		for index, value := range trans {
			fmt.Printf("trans[%d]=%s \n", index, value)
		}

		//fmt.Println(req2map)
	} else {
		fmt.Println(err)
	}

}
func getTransactions(txID string) (str *protos.Transaction, err error) {
	var URL = loc_server + "/transactions/" + txID
	var b []byte
	if response, err := http.Get(URL); err != nil {
		return str, err
	} else {
		if b, err = ioutil.ReadAll(response.Body); err != nil {
			return str, err
		}
	}
	var t0 protos.Transaction
	err = json.Unmarshal(b, &t0)
	if err == nil {
		return &t0, nil
	} else {
		fmt.Println(err)
		return nil, err
	}

}
