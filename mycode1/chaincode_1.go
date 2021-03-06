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
	"crypto/x509"
	"github.com/hyperledger/fabric/core/crypto/primitives"
)
	//"crypto/md5"
	//"crypto/rand"
	//"encoding/base64"
	//"encoding/hex"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
var log=shim.NewLogger("mycode1");

type Account struct{
	AccountNo string
	CustName string
	Amount float64
}
// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	//log.SetLevel(shim.LogDebug)
	log.Infof("main:")
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
		log.Errorf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//if len(args) != 1 {
	//	return nil, errors.New("Incorrect number of arguments. Expecting 1")
	//}
	//log.SetLevel(shim.LogDebug)

	log.Warningf("function:"+function)
	if function == "init"{
		return t.createAccount(stub,args)
	}else if function == "createAccount"{
		return t.createAccount(stub,args)
	}
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)
	log.Warningf("invoke function:"+function)

	caller,err:=stub.GetCallerMetadata()
	if err != nil {
		log.Warningf("Invoke GetCallerMetadata ERR: [%s]" , err.Error())
	}
	log.Infof("Invoke GetCallerMetadata: [%x][%s]" , caller,caller)
	var tcert *x509.Certificate
	certRaw,err:=stub.GetCallerCertificate()
	if err != nil {
		log.Warningf("Invoke GetCallerCertificate ERR: [%s]" , err.Error())
	}
	//log.Infof("Invoke GetCallerCertificate: [%x][%s]" , caller,caller)

	tcert, err = primitives.DERToX509Certificate(certRaw)
	if err != nil {
		log.Warningf("Invoke DERToX509Certificate ERR: [%s]" , err.Error())
	}
	log.Infof("Invoke GetCallerCertificate tcert: [%s][%s]" , tcert.Subject,tcert)
	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}else if function == "createAccount"{
		return t.createAccount(stub,args)
	}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	log.Infof("query function:"+function)
	// Handle different functions
	if function == "dummy_query" {											//read a variable
		fmt.Println("hi there " + function)						//error
		return nil, nil;
	}else if function == "getAccount"{
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		log.Infof("query args:"+args[0])
		_,actBytes, err := t.getAccount(stub,args[0])
		if err != nil {
			fmt.Println("error get Account")
			return nil, err
		}
		return actBytes, nil
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

	log.Infof("args[0]: " + args[0])
	log.Infof("args[1]: " + args[1])
	log.Infof("args[2]: " + args[2])

	f, err := strconv.ParseFloat(args[2], 32)

	account = Account {AccountNo:args[0],CustName:args[1],Amount:f}
	err = t.writeAccount(stub,account)
	if err != nil{
		return nil, errors.New("write Error" + err.Error())
	}

	accountBytes ,err = json.Marshal(&account)
	if err!= nil{
		return nil,errors.New("Error retrieving schoolBytes")
	}


	return accountBytes,nil
}

func (t *SimpleChaincode) writeAccount(stub shim.ChaincodeStubInterface,account Account)(error){
	actBytes ,err := json.Marshal(&account)
	if err != nil{
		return err
	}

	//stub.RangeQueryState("","")
	log.Infof("writeAccount actBytes: %s", actBytes)
	log.Infof("writeAccount accountNo: " + account.AccountNo)
	err = stub.PutState(account.AccountNo,actBytes)
	if err !=nil{
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func (t *SimpleChaincode) getAccount(stub shim.ChaincodeStubInterface,accountNo string)(Account,[]byte,error){
	var account Account
	log.Infof("getAccount accountNo: %s", accountNo)
	actBytes,err := stub.GetState(accountNo)
	if err != nil{
		fmt.Println("Error retrieving data")
	}
	log.Infof("getAccount actBytes: %s", actBytes)
	err = json.Unmarshal(actBytes,&account)
	if err != nil{
		fmt.Println("Error unmarshalling data")
	}
	log.Infof("getAccount accountNo: %s" , account.AccountNo)
	return account,actBytes,nil
}