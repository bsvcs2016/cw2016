/*
Chain code for Testing Blockchain
*/

package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"encoding/json"
	//"strconv"
	//"crypto/x509"
	"strconv"
	//"time"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Issuer struct{
 IssuerID string
 IssuerName string
 IssuerAddress string
 IssuerCreditRating string
 InstrumentList []string
}

type Bank struct {
BankID string
BankName string
BankAddress string
BankTier string
}

type Investor struct {
InvestorID string
InvestorName string
InvestorAddress string
InvestorCreditRating string
}

type Regulator struct {
}

type Transaction struct {
TradeID  string
InstrumentID string
BankID string
}

type Instrument struct {
InstrumentID string
Type string
Maturity string //time.Time
Coupon float64
Notional float64
IssueDate string //time.Time
IssuerID string
BookedNotional float64 
}


func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	issuer:= Issuer{		 
	IssuerID: "Test1",	   
	IssuerName:	"Issuer 1", 
	IssuerAddress: "Client", 
	IssuerCreditRating: "A", 
} 
b, err := json.Marshal(issuer) 
if err == nil { 
       err = stub.PutState(issuer.IssuerID,b) 
}

	investor:= Investor{		 
	InvestorID: "Inv1",	   
	InvestorName:	"Investor 1", 
	InvestorAddress: "Investor NYC", 
	InvestorCreditRating: "A+", 
} 
c, err := json.Marshal(investor) 

if err == nil { 
       err = stub.PutState(investor.InvestorID,c) 
}

bank:= Bank{		 
	BankID: "BOFA",	   
	BankName:	"Bank of America", 
	BankAddress: "NYC", 
	BankTier: "A", 
} 
n, err := json.Marshal(bank) 
if err == nil { 
       err = stub.PutState(bank.BankID,n) 
}
	return nil, nil
}

// Invoke entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "createIssue" {
		return t.createIssue(stub, args)
	} else if function == "invest" {
		return t.write(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) createIssue(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var transactionID string
	var err error
	var vtype, vMaturity,  vIssueDate, vIssuerID string
	var vNotional , vCoupon float64
	fmt.Println("running createIssue")

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5. to create Bond ")
	}
	
	ctidByte, err := stub.GetState("currentTransactionNum")
	if(err != nil){ 
		return nil, errors.New("Error while getting currentTransactionNum from ledger") 
 	} 
	tid,err := strconv.Atoi(string(ctidByte)) 
	if(err != nil){ 
		return nil, errors.New("Error while converting ctidByte to integer") 
	} 

	tid = tid + 1 
	transactionID = "trans"+strconv.Itoa(tid) 

	vtype = args[0] //type
	vMaturity = args[1] 
	vCoupon, err1 := strconv.ParseFloat(args[2],64)
	if err1 != nil {
	return nil, errors.New("Invalid number in Argument 3") 
	}
	vNotional, err2  := strconv.ParseFloat(args[3],64)
	if err2 != nil {
	return nil, errors.New("Invalid number in Argument 4") 
	}
	vIssueDate = args[4] 
	vIssuerID = "Test"  // TODO pull this from User login table
	
	inst := Instrument{
	InstrumentID:transactionID,
	Type:vtype,
	Maturity : vMaturity,
	Coupon:vCoupon,
	Notional:vNotional,
	IssueDate:vIssueDate,
	IssuerID:vIssuerID,
	}
	
	b, err := json.Marshal(inst)
	if err == nil {
    	err = stub.PutState(transactionID, []byte(b)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	}
	err = updateIssueHistory(stub,vIssuerID, transactionID)
	if err != nil {
		errors.New("Error while writing Updating Issuer History into ledger")
	}
	return nil, nil
}
// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func updateIssueHistory(stub shim.ChaincodeStubInterface, issuerID string, instrumentID string) (error) {
	// read entity state
	issuerbyte,err := stub.GetState(issuerID)																										
	if err != nil {
		return errors.New("Error while getting Issuer info from ledger")
	}
	var issuer Issuer
	err = json.Unmarshal(issuerbyte, &issuer)		
	if err != nil {
		return errors.New("Error while unmarshalling issuer data")
	}
	// add InstrumentID to history
	issuer.InstrumentList = append(issuer.InstrumentList,instrumentID)
	// write Issuer state to ledger
	b, err := json.Marshal(issuer)
	if err == nil {
		err = stub.PutState(issuer.IssuerID,b)
	} else {
		return errors.New("Error while updating Issuer status")
	}
	return nil
}
