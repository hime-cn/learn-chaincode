/*
hime test code1
*/

package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"encoding/json"
)
	//"crypto/md5"
	//"crypto/rand"
	//"encoding/base64"
	//"encoding/hex"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Account struct{
	accountNo string
	custName string
	amount float64
}
// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//if len(args) != 1 {
	//	return nil, errors.New("Incorrect number of arguments. Expecting 1")
	//}
	if function == "createAccount"{
		return t.createAccount(stub,args)
	}
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}
	if function == "createAccount"{
		return t.createAccount(stub,args)
	}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "dummy_query" {											//read a variable
		fmt.Println("hi there " + function)						//error
		return nil, nil;
	}else if function == "getAccount"{
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		_,schBytes, err := getAccount(stub,args[0])
		if err != nil {
			fmt.Println("error get Account")
			return nil, err
		}
		return schBytes, nil
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) createAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//if len(args) != 2{
	//	return nil, errors.New("Incorrect number of arguments. Expecting 2")
	//}
	var account Account
	var accountBytes []byte
	//var stuAddress []string
	//var address,priKey,pubKey string
	//address,priKey,pubKey = GetAddress()
	f, err := strconv.ParseFloat(args[2], 32)
	fmt.Println("args[0]: " + args[0])
	fmt.Println("args[1]: " + args[1])
	fmt.Println("args[2]: " + args[2])
	account = Account {accountNo:args[0],custName:args[1],amount:f}
	err = writeAccount(stub,account)
	if err != nil{
		return nil, errors.New("write Error" + err.Error())
	}

	accountBytes ,err = json.Marshal(&account)
	if err!= nil{
		return nil,errors.New("Error retrieving schoolBytes")
	}


	return accountBytes,nil
}

func writeAccount(stub shim.ChaincodeStubInterface,account Account)(error){
	actBytes ,err := json.Marshal(&account)
	if err != nil{
		return err
	}
	//stub.RangeQueryState("","")
	err = stub.PutState(account.accountNo,actBytes)
	if err !=nil{
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func getAccount(stub shim.ChaincodeStubInterface,accountNo string)(Account,[]byte,error){
	var account Account
	actBytes,err := stub.GetState(accountNo)
	if err != nil{
		fmt.Println("Error retrieving data")
	}

	err = json.Unmarshal(actBytes,&account)
	if err != nil{
		fmt.Println("Error unmarshalling data")
	}
	return account,actBytes,nil
}