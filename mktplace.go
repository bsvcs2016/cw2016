/*
Chain code for Test
*/

package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"encoding/json"
	"encoding/pem"
	"strings"
	"crypto/x509"
	"strconv"
	//"time"
	"net/url"
)

//==============================================================================================================================
//	 Participant types - Each participant type is mapped to an integer which we use to compare to the value stored in a
//						 user's eCert
//==============================================================================================================================
//CURRENT WORKAROUND USES ROLES CHANGE WHEN OWN USERS CAN BE CREATED SO THAT IT READ 1, 2, 3, 4
const   ISSUER  =  1
const   BANK   =  2
const   INVESTOR =  3
const   REGBODY  =  4


//==============================================================================================================================
//	 Status types - Asset lifecycle is broken down into 5 statuses, this is part of the business logic to determine what can 
//					be done to the vehicle at points in it's lifecycle
//==============================================================================================================================
const   STATE_INITIAL  			=  0
const   STATE_PARTIAL  			=  1
const   STATE_FULL 		=  2
const   STATE_PAID_OUT 			=  3


// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Issuer struct{
 IssuerID string
 IssuerName string
 IssuerAddress string
 IssuerCreditRating string
 InstrumentList []string
 TradeList []string
}

type Bank struct {
BankID string
BankName string
BankAddress string
BankTier string
InstrumentList []string
TradeList []string
}

type Investor struct {
InvestorID string
InvestorName string
InvestorAddress string
InvestorCreditRating string
InstrumentList []string
TradeList []string
}

type Regulator struct {
}

type Transaction struct {
TradeID  string
InstrumentID string
BankID string
IssuerID string
InvestorID string
Notional float64
Rate float64
}

type Instrument struct {
InstrumentID string
//Type string
Coupon string
Units int
Notional float64
IssueDate string //time.Time
Maturity string //time.Time
IssuerID string
Callable string
Description string
}

/*type Couponpay struct{
}

type Maturitypay struct{
}
*/

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Initialize Two Issuer
	issuer:= Issuer{		 
	IssuerID: "ISS-1",	   
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

	//caller, caller_affiliation, err := t.get_caller_data(stub)
	
	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "createIssue" {
		return t.createIssue(stub, args)
	} else if function == "bankResponse" {
		return t.write(stub, args)
	} else if function == "investorResponse" {
		return t.write(stub, args)		
	}
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
	var vMaturity,  vIssueDate, vIssuerID, vCoupon string
	var vNotional  float64
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

	//vtype = args[0] //type
	vMaturity = args[1] 
	vCoupon  = args[2]
	
	vNotional, err2  := strconv.ParseFloat(args[3],64)
	if err2 != nil {
	return nil, errors.New("Invalid number in Argument 4") 
	}
	vIssueDate = args[4] 
	vIssuerID = "Test"  // TODO pull this from User login table
	
	inst := Instrument{
	InstrumentID:transactionID,
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


//==============================================================================================================================
//	 General Functions
//==============================================================================================================================
//	 get_ecert - Takes the name passed and calls out to the REST API for HyperLedger to retrieve the ecert
//				 for that user. Returns the ecert as retrived including html encoding.
//==============================================================================================================================
func (t *SimpleChaincode) get_ecert(stub shim.ChaincodeStubInterface, name string) ([]byte, error) {
	
	ecert, err := stub.GetState(name)

	if err != nil { return nil, errors.New("Couldn't retrieve ecert for user " + name) }
	
	return ecert, nil
}

//==============================================================================================================================
//	 add_ecert - Adds a new ecert and user pair to the table of ecerts
//==============================================================================================================================

func (t *SimpleChaincode) add_ecert(stub shim.ChaincodeStubInterface, name string, ecert string) ([]byte, error) {
	
	
	err := stub.PutState(name, []byte(ecert))

	if err == nil {
		return nil, errors.New("Error storing eCert for user " + name + " identity: " + ecert)
	}
	
	return nil, nil

}
//==============================================================================================================================
//	 get_caller - Retrieves the username of the user who invoked the chaincode.
//				  Returns the username as a string.
//==============================================================================================================================

func (t *SimpleChaincode) get_username(stub shim.ChaincodeStubInterface) (string, error) {

	bytes, err := stub.GetCallerCertificate();
	if err != nil { return "", errors.New("Couldn't retrieve caller certificate") }
	x509Cert, err := x509.ParseCertificate(bytes);				// Extract Certificate from result of GetCallerCertificate						
															if err != nil { return "", errors.New("Couldn't parse certificate")	}
															
	return x509Cert.Subject.CommonName, nil
}


//==============================================================================================================================
//	 get_caller_data - Calls the get_ecert and check_role functions and returns the ecert and role for the
//					 name passed.
//==============================================================================================================================

func (t *SimpleChaincode) get_caller_data(stub shim.ChaincodeStubInterface) (string, int, error){	

	user, err := t.get_username(stub)
		if err != nil { return "", -1, err }

	ecert, err := t.get_ecert(stub, user);					
		if err != nil { return "", -1, err }

	affiliation, err := t.check_affiliation(stub,string(ecert));			
		if err != nil { return "", -1, err }

	return user, affiliation, nil
}


//==============================================================================================================================
//	 check_affiliation - Takes an ecert as a string, decodes it to remove html encoding then parses it and checks the
// 				  		certificates common name. The affiliation is stored as part of the common name.
//==============================================================================================================================

func (t *SimpleChaincode) check_affiliation(stub shim.ChaincodeStubInterface, cert string) (int, error) {																																																					
	

	decodedCert, err := url.QueryUnescape(cert);    				// make % etc normal //
	
	if err != nil { return -1, errors.New("Could not decode certificate") }
	
	pem, _ := pem.Decode([]byte(decodedCert))           				// Make Plain text   //

	x509Cert, err := x509.ParseCertificate(pem.Bytes);				// Extract Certificate from argument //
														
	if err != nil { return -1, errors.New("Couldn't parse certificate")	}

	cn := x509Cert.Subject.CommonName
	
	res := strings.Split(cn,"\\")
	
	affiliation, _ := strconv.Atoi(res[2])
	
	return affiliation, nil
		
}
