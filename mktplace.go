package main
import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"encoding/json"
	"strconv"
	"crypto/x509"
	"strings"
	"time"
	//"encoding/pem"
	//"net/url"
	
)
type Stock struct{
	Symbol string
	Client string
	Quantity int
	Commission float64
}
type Instrument struct{
	Symbol string
	Coupon string
	Quantity int
	InstrumentPrice float64
	Rate float64
	SettlementDate string
	IssueDate	string
	Callable	string
	TradeID []string
	QuantityResponded int
	Status string
	Owner string
	Bank string
	Issuer string
}
type Entity struct{
	EntityID string				// enrollmentID
	EntityName string
	EntityType string
	Portfolio []Stock
	Instruments []string
	TradeHistory []string		// list of tradeIDs
	IoiList []string
	Balance float64
}

type Transaction struct{		// ledger transactions
	TransactionID string		// different for every transaction
	TradeID string				// same for all transactions corresponding to a single trade
	TransactionType string		// type of transaction rfq or resp or tradeExec or tradeSet	   Request	Response Execute	Exercise
	FromUser string				// entityId of client
	ToUser string				// entityId of bank1 or bank2
	Symbol string				
	Quantity int
	InstrumentPrice float64
	Rate float64	
	SettlementDate time.Time	
	Status string
	TimeStamp string
}

type Trade struct				
{
	TradeID string				// rfq transaction id
	Symbol string
	Quantity int
	TradeType string			// Not Required
	TransactionHistory []string // transactions belonging to this trade
	Status string				// "New Issue" or "Bank Response" or "Issuer Accepted" or "Pending Allocation" or "Trade timed out"
}

type Ioi struct				
{
	IoiId string				// ioi/rfq id
	Notional float64
	Tenor string
	Bank string
	TransactionHistory []string // transactions belonging to this Ioi
	Status string				// "New" or "Responded" 
	Symbol string
	Owner  string
}

const entity1 = "user_type1_1" //issuer1
const entity2 = "user_type1_2"  //issuer2
const entity3 = "user_type1_3"  //bank1
const entity4 = "user_type1_4"	//bank2
const entity7 = "user_type2_2"	//Investor1
const entity8 = "user_type2_3"  //investor2

const entity9 = "user_type2_4"  //regulator

type SimpleChaincode struct {
}
func main() {
    err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Printf("Error starting chaincode: %s", err)
    }
}
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// initialize Instruent	
	/*
	instrument:= Instrument{		
		Symbol :"ISU90D10BPS",
	Coupon :"BW",
	Quantity :10000,
	InstrumentPrice :100,
	Rate :.1,
	SettlementDate :"03/11/2017",
	IssueDate	:"12/14/2016",
	Callable	:"Yes",
	Status :"New Issue",
	Owner :"user_type1_1",
	Issuer :"user_type1_1",
	}
	b, err := json.Marshal(instrument)
	if err == nil {
        err = stub.PutState(instrument.Symbol,b)
    } else {
		return nil, err
	}
	
	instrument1:= Instrument{		
		Symbol :"ISU120DCALL",
	Coupon :"M",
	Quantity :50000,
	InstrumentPrice :50,
	Rate :.2,
	SettlementDate :"07/11/2017",
	IssueDate	:"01/11/2017",
	Callable	:"Yes",
	Status :"New Issue",
	Owner :"user_type1_1",
	Issuer :"user_type1_1",
	}
	b, err = json.Marshal(instrument1)
	if err == nil {
        err = stub.PutState(instrument1.Symbol,b)
    } else {
		return nil, err
	}
	*/
	// initialize entities	
	client:= Entity{		
		EntityID: entity1,	  
		EntityName:	"ZocDoc, Inc.",
		EntityType: "Issuer",
		Balance : 50000000.00,
	}
	//client.Instruments = append(client.Instruments,"ISU90D10BPS")
	//client.Instruments = append(client.Instruments,"ISU120DCALL")
	b, err := json.Marshal(client)
	if err == nil {
        err = stub.PutState(client.EntityID,b)
    } else {
		return nil, err
	}
	
	client2:= Entity{		
		EntityID: entity2,	  
		EntityName:	"Uber Technologies Inc.",
		EntityType: "Issuer",
		Balance : 50000000.00,
	}
	b1, err := json.Marshal(client2)
	if err == nil {
        err = stub.PutState(client2.EntityID,b1)
    } else {
		return nil, err
	}
	
	bank1:= Entity{
		EntityID: entity3,
		EntityName:	"Bank of America Corporation",
		EntityType: "Bank",
		Balance : 100000000.00,
	}
	b, err = json.Marshal(bank1)
	if err == nil {
        err = stub.PutState(bank1.EntityID,b)
    } else {
		return nil, err
	}
	bank2:= Entity{
		EntityID: entity4,
		EntityName:	"Barclays PLC",
		EntityType: "Bank",
		Balance : 100000000.00,
	}
	b, err = json.Marshal(bank2)
	if err == nil {
		err = stub.PutState(bank2.EntityID,b)
    } else {
		return nil, err
	}
	
	regBody:= Entity{
		EntityID: entity9,
		EntityName:	"U.S. Securities and Exchange Commission (SEC)",
		EntityType: "RegBody",
	}
	b, err = json.Marshal(regBody)
	if err == nil {
		err = stub.PutState(regBody.EntityID,b)
    } else {
		return nil, err
	}
	
	inv1:= Entity{
		EntityID: entity7,
		EntityName:	"Fidelity Investments Inc.",
		EntityType: "Investor",
		Balance : 200000000.00,
	}
	b, err = json.Marshal(inv1)
	if err == nil {
		err = stub.PutState(inv1.EntityID,b)
    } else {
		return nil, err
	}
	
	inv2:= Entity{
		EntityID: entity8,
		EntityName:	"BlackRock, Inc.",
		EntityType: "Investor",
		Balance : 200000000.00,
	}
	b, err = json.Marshal(inv2)
	if err == nil {
		err = stub.PutState(inv2.EntityID,b)
    } else {
		return nil, err
	}
	
	EntityList := []string{entity1,entity2, entity3, entity4, entity7, entity8, entity9}

	b, err = json.Marshal(EntityList)
	if err == nil {
		err = stub.PutState("entityList",b)
    } else {
		return nil, err
	}
	
	// initialize trade num and transaction num
	byteVal, err := stub.GetState("currentTransactionNum")
	if len(byteVal) == 0 {
		err = stub.PutState("currentTransactionNum", []byte("1000"))
	}
	ctidByte,err := stub.GetState("currentTransactionNum")
	if(err != nil){
		return nil, errors.New("Error while getting currentTransactionNum from ledger")
	}
	
	byteVal, err = stub.GetState("currentTradeNum")
	if len(byteVal) == 0 {
		err = stub.PutState("currentTradeNum", []byte("1000"))
	}
	ctidByte,err = stub.GetState("currentTradeNum")
	if(err != nil){
		return nil, errors.New("Error while getting currentTradeNum from ledger")
	}
	
	// initialize Ioi num
	byteVal, err = stub.GetState("currentIoiNum")
	if len(byteVal) == 0 {
		err = stub.PutState("currentIoiNum", []byte("1000"))
	}
	ctidByte,err = stub.GetState("currentIoiNum")
	if(err != nil){
		return nil, errors.New("Error while getting currentIoiNum from ledger")
	}
	// initialize Instrument num 
	byteVal, err = stub.GetState("currentInstrumentNum")
	if len(byteVal) == 0 {
		err = stub.PutState("currentInstrumentNum", []byte("1000"))
	}
	ctidByte,err = stub.GetState("currentInstrumentNum")
	if(err != nil){
		return nil, errors.New("Error while getting currentInstrumentNum from ledger")
	}
    return ctidByte, nil
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

    username, err := stub.ReadCertAttribute("enrollmentID");
	if err != nil { return "", errors.New("Couldn't get attribute 'username'. Error: " + err.Error()) }
	return string(username), nil
}
//==============================================================================================================================
//	 check_affiliation - Takes an ecert as a string, decodes it to remove html encoding then parses it and checks the
// 				  		certificates common name. The affiliation is stored as part of the common name.
//==============================================================================================================================

func (t *SimpleChaincode) check_affiliation(stub shim.ChaincodeStubInterface) (string, error) {
    affiliation, err := stub.ReadCertAttribute("enrollmentID");
	if err != nil { return "", errors.New("Couldn't get attribute 'role'. Error: " + err.Error()) }
	return string(affiliation), nil

}


//==============================================================================================================================
//	 get_caller_data - Calls the get_ecert and check_role functions and returns the ecert and role for the
//					 name passed.
//==============================================================================================================================

func (t *SimpleChaincode) get_caller_data(stub shim.ChaincodeStubInterface) (string, string, error){

	user, err := t.get_username(stub)

    // if err != nil { return "", "", err }

	// ecert, err := t.get_ecert(stub, user);

    // if err != nil { return "", "", err }

	affiliation, err := t.check_affiliation(stub);

    if err != nil { return "", "", err }

	return user, affiliation, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    
	caller, err := t.get_username(stub)
	fmt.Println("Caller Detail " + caller+" :"  +err.Error())
	
	// Handle different functions
    if function == "init" {
        return t.Init(stub, "init", args)
    } else if function == "createIssue" {
        return t.createIssue(stub, args)
    } else if function == "requestForIssue" {
        return t.requestForIssue(stub, args)
    } else if function == "respondToIssue" { //Pass Response as well (Bank/Investor)
        return t.respondToIssue(stub, args)
    } else if function == "tradeExec" {  
        return t.tradeExec(stub, args)
	} else if function == "tradeSet" {     // Money and Coupon price will be transfered to Bank and From Bank to Investors
        return t.tradeSet(stub, args)
    } else if function == "trial" {
        return t.trial(stub, args)
	} else if function == "payCoupon" {
        return t.payCoupon(stub, args)
	} else if function == "issueCallout" {
        return t.issueCallout(stub, args)
	} else if function == "requestForInstrument" {
        return t.requestForInstrument(stub, args)
    } 
    fmt.Println("invoke did not find func: " + function)
    return nil, errors.New("Received unknown function invocation")
}
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    // Handle different functions
    if function == "readEntity" {
        return t.readEntity(stub, args)
    }	else if function =="readTransaction" {
		return t.readTransaction(stub,args)
	}	else if function =="getUserID" {
		return t.getUserID(stub,args)
	}	else if function =="getcurrentTransactionNum" {
		return t.getcurrentTransactionNum(stub,args)
	}	else if function == "getValue" {
        return t.getValue(stub, args)
	}	else if function == "readTradeIDsOfUser" {
        return t.readTradeIDsOfUser(stub, args)
    }	else if function == "readTrades" {
        return t.readTrades(stub, args)
    }	else if function == "readIssueRequests" {
        return t.readIssueRequests(stub, args)
    }	else if function == "getAllTrades" {
        return t.getAllTrades(stub, args)
    }	else if function == "getEntityList" {
        return t.getEntityList(stub, args)
	}	else if function == "getEntities" {
        return t.getEntities(stub, args)
    }	else if function == "getTransactionStatus" {
        return t.getTransactionStatus(stub, args)
	}	else if function == "getInstrument" {
        return t.getInstrument(stub, args)
	}	else if function == "getAllInstruments" {
        return t.getAllInstruments(stub, args)
	}	else if function == "getAllInstrumentTrades" {
        return t.getAllInstrumentTrades(stub, args)
	}	else if function == "getAllIoi" {
        return t.getAllIoi(stub, args)
    }
	fmt.Println("query did not find func: " + function)
    return nil, errors.New("Received unknown function query")
}
func (t *SimpleChaincode) readEntity(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var jsonResp string
    var err error
	var valAsbytes []byte
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting entity ID")
    }
	valAsbytes, err = stub.GetState(args[0])
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
        return nil, errors.New(jsonResp)
    }
    return valAsbytes, nil
}
func (t *SimpleChaincode) readTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var tid, jsonResp string
    var err error
    if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting transaction ID")
    }
    tid = args[0]
    valAsbytes, err := stub.GetState(tid)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + tid + "\"}"
		return nil, errors.New(jsonResp)
    }
	var tran Transaction
	err = json.Unmarshal(valAsbytes, &tran)
	if(err != nil){
		return nil, errors.New("Error while unmarshalling transaction data")
	}
	return valAsbytes, nil
	bytes, err := stub.GetCallerCertificate();
	if(err != nil){
		return nil, errors.New("Error while getting caller certificate")
	}
	x509Cert, err := x509.ParseCertificate(bytes);
	
	// check entity type and accordingly allow transaction to be read
	entityByte,err := stub.GetState(x509Cert.Subject.CommonName)
	if(err != nil){
		return nil, errors.New("Error while getting bank info from ledger")
	}
	var entity Entity
	err = json.Unmarshal(entityByte, &entity)
	if(err != nil){
		return nil, errors.New("Error while unmarshalling entity data")
	}
	
	switch entity.EntityType {
		case "RegBody":	return valAsbytes, nil
		case "Issuer":	if tran.FromUser == x509Cert.Subject.CommonName {
							return valAsbytes, nil
						}
		case "Bank":	if tran.ToUser == x509Cert.Subject.CommonName {
							return valAsbytes, nil
						}
		case "Investor": if tran.FromUser == x509Cert.Subject.CommonName {
						 return valAsbytes, nil
						}
	}
    return nil, nil
}
// used by Client to send to Banks for new Issue.
/*		arg 0 	: caller
		arg 1	:	Symbol
		arg 2 	: Bank/Investor
		arg 3	: Status /PublishToBank/PublishToInvestor
		
	
*/
func (t *SimpleChaincode) requestForIssue(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//Need all parameters for the Bond Instrument
	if len(args) >=3{
		var transactionID  string 
		//get instrument detail
		instbyte , err := stub.GetState(args[1])
		var instr Instrument
		err = json.Unmarshal(instbyte, &instr)
		if(err != nil){
			return nil, errors.New("Error while unmarshalling Instrument data:" +args[1])
		}
		fmt.Println("Owner "+instr.Owner)
		fmt.Println("Caller "+args[0])
		/*
		if instr.Owner != args[0] {    // if caller is not the owner then return error 
			return nil, errors.New("Only Owner can Request for Investment")
		}
		*/
		//for i :=2; i < len(args); i++ {
		// get current Trade number
		// get current Transaction number
		
		ctidByte1,err1 := stub.GetState("currentTransactionNum")
		if(err1 != nil){
			return nil, errors.New("Error while getting currentTransactionNum from ledger")
		}
		tid,err := strconv.Atoi(string(ctidByte1))
		if(err != nil){
			return nil, errors.New("Error while converting ctidByte to integer")
		}
		tid = tid + 1
		transactionID = "trans"+strconv.Itoa(tid)
		
		if(err != nil){
			_ = updateTransactionStatus(stub, transactionID, "Error while converting ctidByte to integer")
			return nil, nil
		}

		bytes, err := stub.GetCallerCertificate();
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while getting caller certificate")
			return nil, nil
		}
		// get client enrollmentID
		x509Cert, err := x509.ParseCertificate(bytes);
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while parsing caller certificate")
			return nil, nil
		}
		fmt.Println("x509Cert.Subject.CommonName :" +x509Cert.Subject.CommonName)
		status := args[3]
		//  Create Multiple Transactions with Each Bank as per selection in UI
		//Transaction
		trn := Transaction{
		TransactionID: transactionID,
		TransactionType: status,
		FromUser:	args[0],	// enrollmentID
		ToUser: args[2],
		Symbol: args[1],						// based on input
		Quantity:	instr.Quantity,								// based on input
		InstrumentPrice: instr.InstrumentPrice,
		Rate: instr.Rate,
		Status: "Success",
		TimeStamp : time.Now().Format("2006-01-02 15:04:05"),
		}
		fmt.Println("Transaction")
		//clientID = trn.FromUser
		// convert to Transaction to JSON
		b, err := json.Marshal(trn)
		// write to ledger
		if err == nil {
			err = stub.PutState(trn.TransactionID,b)
			if err != nil {
				_ = updateTransactionStatus(stub, transactionID, "Error while writing Transaction to ledger")
				return nil, nil
			}
		} else {
			_ = updateTransactionStatus(stub, transactionID, "Error while marshalling trade data")
			return nil, nil
		}
		fmt.Println("Transaction Updated")
		bankByte, err := stub.GetState(args[2])
		if err != nil {
			return nil, errors.New("Unable to get Bank's data")
		}
		var bank Entity
		err = json.Unmarshal(bankByte,&bank)
		if err != nil {
			return nil, errors.New("Unable to unmarshal Bank's data")
		}
		
		// update currentTransactionNum
		err = stub.PutState("currentTransactionNum", []byte(strconv.Itoa(tid)))
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while updating current transaction number")
			return nil, nil
		}
		
		// add Transaction ID to entity's trade history
		err = updateTradeHistory(stub, trn.ToUser, trn.TransactionID)
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while updating trade history")
			return nil, nil
		}	
		// add Transaction ID to entity's trade history
		err = updateTradeHistory(stub, trn.FromUser, trn.TransactionID)
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while updating trade history for Issuer")
			return nil, nil
		}
		fmt.Println("Instruent History" +trn.FromUser)
		err = updateInstrumentHistory(stub, trn.ToUser, trn.Symbol)
		if err != nil {
			return nil, errors.New( "Error while updating Instrument History : Caller : "+trn.ToUser+" :"+trn.Symbol)
		}
		fmt.Println("Instruent History" +trn.FromUser)
		err = t.updateInstrumentStatus(stub, trn.Symbol, args[2], status)
		if err != nil {
			return nil, errors.New( "Error while updating Instrument History : Caller : "+trn.ToUser+" :"+trn.Symbol)
		}
		fmt.Println("Instruent updateInstrumentStatus" +trn.TransactionID)
		err = t.updateInstrumentTradeHistory(stub, trn.Symbol, trn.TransactionID)
		if err != nil {
			return nil, errors.New( "Error while updating Instrument Trade Histiry History : Caller : "+trn.TransactionID+" :"+trn.Symbol)
		}
		
	 //} //For loop
		return []byte(transactionID), nil
	}
	return nil, errors.New("Incorrect number of arguments")
}
/*			arg 0	:	Caller
			arg 1	:	Instrument ID
			arg 2	:	Response (yes/no)
			arg 3	:	Status
*/
func (t *SimpleChaincode) respondToIssue(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args)== 4 {
		caller := args[0]
		symbol := args[1]
		response := args[2]
		status := args[3]
		instbyte, err := stub.GetState(symbol)
		if err != nil {
			return nil, errors.New("Instruent not found")
		}
		var inst Instrument
		err= json.Unmarshal(instbyte, &inst)
		if err != nil {
			return nil, errors.New("Error in unmarshalling instruent ")
		}
		quoteID := inst.TradeID[len(inst.TradeID)-1]
		
		// get information from requestForIssue transaction
		rfqbyte,err := stub.GetState(quoteID)												
		if err != nil {
			_ = updateTransactionStatus(stub, quoteID, "Error while reading quote request transaction from ledger")
			return nil, nil
		}
		var rfq Transaction
		err = json.Unmarshal(rfqbyte, &rfq)
		if err != nil {
			_ = updateTransactionStatus(stub, quoteID, "Error while unmarshalling quote request data")
			return nil, nil
		}

		if response =="yes" {
		ctidByte, err := stub.GetState("currentTransactionNum")
		if(err != nil){
			return nil, errors.New("Error while getting currentTransactionNum from ledger")
		}
		tid,err := strconv.Atoi(string(ctidByte))
		if(err != nil){
			return nil, errors.New("Error while converting ctidByte to integer")
		}
		tid = tid + 1
		transactionID := "trans"+strconv.Itoa(tid)
		
		// get bank's enrollment id
		bytes, err := stub.GetCallerCertificate();
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while getting caller certificate")
			return nil, nil
		}
		x509Cert, err := x509.ParseCertificate(bytes);
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while parsing caller certificate")
			return nil, nil
		}		
		fmt.Println("Respond to Issue : x509Cert"+x509Cert.Subject.CommonName)
		caller1, err := t.get_username(stub)
		fmt.Println("Respond to Issue : Caller1 :"+ caller1)
		
		
		if rfq.Symbol != symbol {
			_ = updateTransactionStatus(stub, transactionID, "Error due to mismatch in tradeIDs")
			return nil, nil
		}		
		fmt.Println("Respond to Issue : Quantity "+args[3])
		
		quantity := rfq.Quantity
		
		// check if required quantity is  under limit
		instrumentByte,err := stub.GetState(rfq.Symbol)																											
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while getting Instrument info from ledger")
			return nil, nil
		}
		fmt.Println("Respond to Issue : Instrument Bytes"+string(instrumentByte))
		var inst Instrument
		err = json.Unmarshal(instrumentByte, &inst)
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while unmarshalling Instrument data")
			return nil, nil
		}
		fmt.Println("Respond to Issue : Instrument symbol"+inst.Symbol)
		

		fmt.Printf("Quantity Instrument :%g" ,inst.Quantity )
		if quantity >inst.Quantity {
		 return nil, errors.New("Response Quantity should be less or equal to requested")
		}
		
			
		entityByte, err := stub.GetState(caller)
		if err != nil {
			return nil,errors.New("Instrument with this ID already Exists, Try a different Name")
			
		}
		var entity Entity
		err = json.Unmarshal(entityByte, &entity)
		if err != nil {
			return nil, errors.New("Entity Not Found")
		}
		
		tr := Transaction {
		TransactionID: transactionID,
		TransactionType: status,
		//InstrumentType: rfq.InstrumentType,														// get from rfq
		FromUser:	caller,														// 
		ToUser: rfq.FromUser,  //x509Cert.Subject.CommonName,											// 
		Symbol: rfq.Symbol,													// get from rfq
		Quantity:	quantity,														// get from rfq
		InstrumentPrice: rfq.InstrumentPrice,																// based on input
		Rate: rfq.Rate,																// based on input
		//SettlementDate: time.Date(year, month, day, 0, 0, 0, 0, time.UTC),				// based on input
		Status: "Success",
		TimeStamp : time.Now().Format("2006-01-02 15:04:05"),
		}

		// convert to JSON
		b, err := json.Marshal(tr)
		
		// write to ledger
		if err == nil {
			err = stub.PutState(tr.TransactionID,b)
			if err != nil {
				_ = updateTransactionStatus(stub, transactionID, "Error while writing Response transaction to ledger")
				return nil, nil
			}
		} else {
			_ = updateTransactionStatus(stub, transactionID, "Error while marshalling transaction data")
			return nil, nil
		}
		
		err = stub.PutState("currentTransactionNum", []byte(strconv.Itoa(tid)))
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while writing current Transaction Number to ledger")
			return nil, nil
		}
		
		inst.QuantityResponded = quantity
		commission := float64(quantity)*inst.InstrumentPrice*.001
		price := float64(quantity)*inst.InstrumentPrice
		err = t.updateEntityBalance(stub,caller, -1*(price - commission))   //Caller
		if err != nil {
				return nil, errors.New(err.Error())
		}
		err = t.updateEntityBalance(stub,rfq.FromUser, (price - commission))   //Client
		if err != nil {
				return nil, errors.New(err.Error())
		}
		b, err = json.Marshal(inst)
		err = stub.PutState(inst.Symbol,b)
		if err != nil{
			return nil, errors.New("Unable to update Instrument Responded Quantity "+err.Error())
		}
		
		err = t.updateInstrumentStatus(stub, args[1],rfq.ToUser,status)
		if err != nil{
		 return nil,errors.New("Unable to update Instruent Status")
		}
		// add Transaction ID to entity's trade history
		err = updateTradeHistory(stub, tr.ToUser, tr.TransactionID)
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while updating trade history")
			return nil, nil
		}	
		// add Transaction ID to entity's trade history
		err = updateTradeHistory(stub, tr.FromUser, tr.TransactionID)
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while updating trade history for Issuer")
			return nil, nil
		}
		fmt.Println("Instruent updateInstrumentStatus" +tr.TransactionID)
		err = t.updateInstrumentTradeHistory(stub, tr.Symbol, tr.TransactionID)
		if err != nil {
			return nil, errors.New( "Error while updating Instrument Trade Histiry History : Caller : "+tr.TransactionID+" :"+tr.Symbol)
		}
		return nil, nil
	}	else{  // not accepted
		err := t.updateInstrumentStatus(stub, args[1],rfq.FromUser,status)
		if err != nil{
		 return nil,errors.New("Unable to update Instruent Status")
		}
	}
	}
	return nil, errors.New("Incorrect number of arguments")
}
/*			arg 0	:	TradeID
			arg 1	:	Selected quote's TransactionID
*/
//---------------------------------------------------------- consensus
func (t *SimpleChaincode) tradeExec(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args)== 4 {
		caller := args[0]
		ctidByte, err := stub.GetState("currentTransactionNum")
		if err != nil {
			return nil, errors.New("Error while getting current Transaction Number from ledger")
		}		
		tid,err := strconv.Atoi(string(ctidByte))
		if err != nil {
			return nil, errors.New("Error while converting ctidByte to integer")
		}
		tid = tid + 1
		transactionID := "trans"+strconv.Itoa(tid)
		
		fmt.Println("Current Transaction No"+transactionID)
		
		tradeID := args[1]
		quoteId := args[2]
		
		// get client's enrollment id
		bytes, err := stub.GetCallerCertificate();
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while getting caller certificate")
			return nil, nil
		}
		x509Cert, err := x509.ParseCertificate(bytes);
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while parsing caller certificate :" + x509Cert.Subject.CommonName)
			return nil, nil
		}
		fmt.Println("Current x509Cert No :"+x509Cert.Subject.CommonName + quoteId)
		// get information from selected quote
		quotebyte,err := stub.GetState(quoteId)
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while getting quote data")
			return nil, nil
		}
		fmt.Println("Quote  Id   :"+string(quotebyte))
		var quote Transaction
		err = json.Unmarshal(quotebyte, &quote)		
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while unmarshalling quote data")
			return nil, nil
		}
		fmt.Println("Trade  ID   :"+quote.TradeID +"-"+ tradeID)
		if quote.TradeID != tradeID {
			_ = updateTransactionStatus(stub, transactionID, "Error due to mismatch in tradeIDs")	
			return nil, nil
		}
		fmt.Println("Quote Trade Id   :"+tradeID)


		// check if trade has to be Executed or Cancelled
		if strings.ToLower(args[3]) == "yes" {
			tExec := quote
			if tExec.TradeID != tradeID {
				_ = updateTransactionStatus(stub, transactionID, "Error due to mismatch in tradeIDs")
				return nil, nil
			}
			
			instByte, err := stub.GetState(tExec.Symbol)
			if err != nil {
				return nil, errors.New("Error pulling Instrument state")
			}
			var inst Instrument
			err = json.Unmarshal(instByte, &inst)
			if err != nil {
				return nil, errors.New("Error unmarshalling Instrument state")
			}

			// check settlement date to see if instrument is still valid
				
				t := Transaction{
				TransactionID: transactionID,
				TradeID: tradeID,							// based on input
				TransactionType: "Final",
				//InstrumentType: tExec.InstrumentType,				// get from tradeExec transaction
				FromUser: caller , //x509Cert.Subject.CommonName,		// get from tradeExec transaction
				ToUser: tExec.ToUser,						// get from tradeExec transaction
				Symbol: tExec.Symbol,				// get from tradeExec transaction
				Quantity:	tExec.Quantity,					// get from tradeExec transaction
				InstrumentPrice: tExec.InstrumentPrice,				// get from tradeExec transaction
				Rate: tExec.Rate,					// get from tradeExec transaction
				Status: "Success",
				}
				// convert to JSON
				b1, err := json.Marshal(t)
				// write to ledger
				if err == nil {
					err = stub.PutState(t.TransactionID,b1)
					if err != nil {
						_ = updateTransactionStatus(stub, transactionID, "Error while writing Response transaction to ledger")
						return nil, nil
					}
				} else {
					_ = updateTransactionStatus(stub, transactionID, "Error while marshalling transaction data")
					return nil, nil
				}
				
						// update client entity's instruments
		clientbyte,err := stub.GetState(t.FromUser)																										
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while getting client info from ledger")
			return nil, nil
		}
		var client Entity
		err = json.Unmarshal(clientbyte, &client)		
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while unmarshalling client data")
			return nil, nil
		}
		
		
		bankbyte,err := stub.GetState(t.ToUser)																										
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while getting bank information from ledger")
			return nil, nil
		}
		var bank Entity
		err = json.Unmarshal(bankbyte, &bank)		
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while unmarshalling bank data")
			return nil, nil
		}
		
		fmt.Println("Bank and Client ID received")

		// add stock to clients portfolio, check if stock already exists if yes increase quantity else create new stock entry 		
		stockExistFlag := false
		for i := 0; i< len(client.Portfolio); i++ {
			if client.Portfolio[i].Symbol == t.Symbol && client.Portfolio[i].Client == t.ToUser {
				stockExistFlag = true
				client.Portfolio[i].Quantity = client.Portfolio[i].Quantity - t.Quantity
				if client.EntityType =="Issuer" {
					client.Portfolio[i].Commission = client.Portfolio[i].Commission + float64(-t.Quantity) * inst.InstrumentPrice *.001
				}else {
					client.Portfolio[i].Commission = client.Portfolio[i].Commission + float64(t.Quantity) * inst.InstrumentPrice *.001
				}
				break
			  }
			}
				// create new stock entry
				if stockExistFlag == false {
					newStock := Stock{Symbol: t.Symbol,Client: t.ToUser, Quantity: t.Quantity, Commission: float64(-t.Quantity) * inst.InstrumentPrice *.001}
					client.Portfolio = append(client.Portfolio,newStock)
				}
				// update banks stock data
				stockExistFlag = false
				for i := 0; i< len(bank.Portfolio); i++ {
					if bank.Portfolio[i].Symbol == t.Symbol  && client.Portfolio[i].Client == t.FromUser {
						stockExistFlag = true
						bank.Portfolio[i].Quantity = bank.Portfolio[i].Quantity + t.Quantity
						if bank.EntityType =="Investor" {
						client.Portfolio[i].Commission = client.Portfolio[i].Commission + float64(-t.Quantity) * inst.InstrumentPrice *.001
						} else {
						client.Portfolio[i].Commission = client.Portfolio[i].Commission + float64(t.Quantity) * inst.InstrumentPrice *.001
						}
						break
					}
				}
				
				// create new stock entry
				newStock := Stock{Symbol: t.Symbol,Client: t.FromUser, Quantity: t.Quantity, Commission : float64(t.Quantity) * inst.InstrumentPrice *.001}
				bank.Portfolio = append(bank.Portfolio,newStock)

				// updating trade state
				err = updateTradeState(stub, t.TradeID, t.TransactionID,"Trade Executed")
				if err != nil {
					_ = updateTransactionStatus(stub, transactionID, "Error while updating trade state")
					return nil, nil
				}
				
		// update client state
		b, err := json.Marshal(client)
		if err == nil {
			err = stub.PutState(client.EntityID,b)
		} else {
			_ = updateTransactionStatus(stub, transactionID, "Error updating Client state")
			return nil, nil
		}
		// update bank state
		b, err = json.Marshal(bank)
		if err == nil {
			err = stub.PutState(bank.EntityID,b)
		} else {
			_ = updateTransactionStatus(stub, transactionID, "Error while updating Bank state")
			return nil, nil
		}
		// update transaction number
		err = stub.PutState("currentTransactionNum", []byte(strconv.Itoa(tid)))
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while writing currentTransactionNum to ledger")
			return nil, nil
		}		

		} else {	// trade cancelled
			_ = updateTransactionStatus(stub, transactionID, "")
			// updating trade state
			err = updateTradeState(stub, tradeID,"" ,"Trade Cancelled")
			if err != nil {
				_ = updateTransactionStatus(stub, transactionID, "Error while updating trade state")
				return nil, nil
			}
		}
	
		return nil, nil
	}
	return nil, errors.New("Incorrect number of arguments")
}
/*			arg 0	:	TradeID
			arg 1	:	Yes/ No
*/

func (t *SimpleChaincode) tradeSet(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
return nil, nil
/*
	if len(args)== 3 {
		caller := args[0]
		tradeID := args[1]
		//tExecId := args[2]
		// get client's enrollment id
		
		ctidByte, err := stub.GetState("currentTransactionNum")
		if(err != nil){
			return nil, errors.New("Error while getting currentTransactionNum from ledger")
		}
		tid,err := strconv.Atoi(string(ctidByte))
		if(err != nil){
			return nil, errors.New("Error while converting ctidByte to integer")
		}	
		tid = tid + 1
		transactionID := "trans"+strconv.Itoa(tid)
		
		bytes, err := stub.GetCallerCertificate();
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while getting caller certificate")
			return nil, nil
		}
		x509Cert, err := x509.ParseCertificate(bytes);
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while parsing caller certificate :" +x509Cert.Subject.CommonName)
			return nil, nil
		}
		clientID := caller //x509Cert.Subject.CommonName
		
		// update client entity's instruments
		clientbyte,err := stub.GetState(clientID)																												
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while getting client info from ledger")
			return nil, nil
		}
		var client Entity
		err = json.Unmarshal(clientbyte, &client)		
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while unmarshalling client data")
			return nil, nil
		}
		// remove instrument from clients data, check tradeID

		// get transactionID from tradeID
		tradebyte,err := stub.GetState(tradeID)
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while getting trade info from ledger")
			return nil, nil
		}
		fmt.Println (" Transaction ID :" + tExecId)
		// get information from trade exec transaction
		tbyte,err := stub.GetState(tExecId)												
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while getting tradeExec transaction from ledger")
			return nil, nil
		}
		
		var tExec Transaction
		err = json.Unmarshal(tbyte, &tExec)		
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while unmarshalling tradeExec data")
			return nil, nil
		}
		
		// update bank entity's instruments
		bankbyte,err := stub.GetState(tExec.ToUser)																											
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while getting bank info from ledger")
			return nil, nil
		}
		var bank Entity
		err = json.Unmarshal(bankbyte, &bank)		
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while unmarshalling bank data")
			return nil, nil
		}
		// remove instrument from bank 

		// check if trade has to be settled
			if tExec.TradeID != tradeID {
				_ = updateTransactionStatus(stub, transactionID, "Error due to mismatch in tradeIDs")
				return nil, nil
			}
			fmt.Println (" Trade ID and Execution Trade Id "+ tExec.TradeID +"-"+ tradeID)
			// check settlement date to see if instrument is still valid
				
				t := Transaction{
				TransactionID: transactionID,
				TradeID: tradeID,							// based on input
				TransactionType: "Final",
				//InstrumentType: tExec.InstrumentType,				// get from tradeExec transaction
				FromUser: caller , //x509Cert.Subject.CommonName,		// get from tradeExec transaction
				ToUser: tExec.ToUser,						// get from tradeExec transaction
				Symbol: tExec.Symbol,				// get from tradeExec transaction
				Quantity:	tExec.Quantity,					// get from tradeExec transaction
				InstrumentPrice: tExec.InstrumentPrice,				// get from tradeExec transaction
				Rate: tExec.Rate,					// get from tradeExec transaction
				Status: "Success",
				}
				// convert to JSON
				b1, err := json.Marshal(t)
				// write to ledger
				if err == nil {
					err = stub.PutState(t.TransactionID,b1)
					if err != nil {
						_ = updateTransactionStatus(stub, transactionID, "Error while writing Response transaction to ledger")
						return nil, nil
					}
				} else {
					_ = updateTransactionStatus(stub, transactionID, "Error while marshalling transaction data")
					return nil, nil
				}
				
				// add stock to clients portfolio, check if stock already exists if yes increase quantity else create new stock entry 		
				//stockExistFlag := false
				for i := 0; i< len(client.Portfolio); i++ {
					if client.Portfolio[i].Symbol == t.Symbol && client.Portfolio[i].Client == t.ToUser {
						//stockExistFlag = true
						client.Portfolio[i].Quantity = client.Portfolio[i].Quantity - t.Quantity
						
						break
					  }
					}
				
				// update banks stock data
				//stockExistFlag = false
				for i := 0; i< len(bank.Portfolio); i++ {
					if bank.Portfolio[i].Symbol == t.Symbol && client.Portfolio[i].Client == t.FromUser{
						//stockExistFlag = true
						bank.Portfolio[i].Quantity = bank.Portfolio[i].Quantity + t.Quantity
						break
					}
				}
				
				// updating trade state
				err = updateTradeState(stub, t.TradeID, t.TransactionID,"Trade Settled")
				if err != nil {
					_ = updateTransactionStatus(stub, transactionID, "Error while updating trade state")
					return nil, nil
				}
				

		// update client state
		b, err := json.Marshal(client)
		if err == nil {
			err = stub.PutState(client.EntityID,b)
		} else {
			_ = updateTransactionStatus(stub, transactionID, "Error updating Client state")
			return nil, nil
		}
		// update bank state
		b, err = json.Marshal(bank)
		if err == nil {
			err = stub.PutState(bank.EntityID,b)
		} else {
			_ = updateTransactionStatus(stub, transactionID, "Error while updating Bank state")
			return nil, nil
		}
		// update transaction number
		err = stub.PutState("currentTransactionNum", []byte(strconv.Itoa(tid)))
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while writing currentTransactionNum to ledger")
			return nil, nil
		}
		return nil, nil
	}
	return nil, errors.New("Incorrect number of arguments")
	*/
}

// get user id
func (t *SimpleChaincode) getUserID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	bytes, err := stub.GetCallerCertificate()
	x509Cert, err := x509.ParseCertificate(bytes)
	return []byte(x509Cert.Subject.CommonName), err
}
func (t *SimpleChaincode) getcurrentTransactionNum(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	ctidByte,err := stub.GetState("currentTransactionNum")
	if err != nil {
		return nil, errors.New("Error retrieving currentTransactionNum")
	}
    return ctidByte, err
}
func (t *SimpleChaincode) getValue(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	byteVal,err := stub.GetState(args[0])
	if err != nil {
		return []byte(err.Error()), errors.New("Error retrieving key "+args[0])
	}
	if len(byteVal) == 0 {
		return []byte("Len is zero"), nil
	}
    return byteVal, nil
}
// read transactions IDs for a particular user
func (t *SimpleChaincode) readTradeIDsOfUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args)== 1 {
		// read entity state
		entitybyte,err := stub.GetState(args[0])																									
		if err != nil {
			return nil, errors.New("Error while getting entity info from ledger")
		}
		var entity Entity
		err = json.Unmarshal(entitybyte, &entity)		
		if(err != nil){
			return nil, errors.New("Error while unmarshalling entity data")
		}

		b, err := json.Marshal(entity.TradeHistory)
		if err != nil {
			return nil, errors.New("Error while marshalling trade history")
		}
		return b, nil
	}
	return nil, errors.New("Incorrect number of arguments")
}
func updateTradeHistory(stub shim.ChaincodeStubInterface, entityID string, tradeID string) (error) {
	// read entity state
	entitybyte,err := stub.GetState(entityID)																										
	if err != nil {
		return errors.New("Error while getting entity info from ledger")
	}
	var entity Entity
	err = json.Unmarshal(entitybyte, &entity)		
	if err != nil {
		return errors.New("Error while unmarshalling entity data")
	}
	// add tradeID to history
	entity.TradeHistory = append(entity.TradeHistory,tradeID)
	// write entity state to ledger
	b, err := json.Marshal(entity)
	if err == nil {
		err = stub.PutState(entity.EntityID,b)
	} else {
		return errors.New("Error while updating entity status")
	}
	return nil
}

func updateInstrumentHistory(stub shim.ChaincodeStubInterface, entityID string, issueID string) (error) {
	// read entity state
	entitybyte,err := stub.GetState(entityID)																										
	if err != nil {
		return errors.New("Error while getting entity info from ledger")
	}
	
	var entity Entity
	err = json.Unmarshal(entitybyte, &entity)		
	if err != nil {
		return errors.New("Error while unmarshalling entity data")
	}
	// add tradeID to history
	//entity.Instruments = append(entity.Instruments,issueID)
	newInstrumeentList := []string{}
	
	for _, v := range entity.Instruments {
			if v != issueID {
		newInstrumeentList = append(newInstrumeentList, v)
			}
		}
	newInstrumeentList = append(newInstrumeentList,issueID)
	entity.Instruments = newInstrumeentList
	
	// write entity state to ledger
	b, err := json.Marshal(entity)
	if err == nil {
		err = stub.PutState(entity.EntityID,b)
	} else {
		return errors.New("Error while updating entity status")
	}
	return nil
}


func updateTradeState(stub shim.ChaincodeStubInterface, tradeID string, transactionID string, status string) (error) {
	return nil
	/*
	// read trade state
	tradebyte,err := stub.GetState(tradeID)																										
	if err != nil {
		return errors.New("Error while getting trade info from ledger")
	}
	var trade Trade
	err = json.Unmarshal(tradebyte, &trade)		
	if err != nil {
		return errors.New("Error while unmarshalling trade data")
	}
	// add transactionID to history
	trade.TransactionHistory = append(trade.TransactionHistory,transactionID)
	
	// update status
	trade.Status = status
	
	// write trade state to ledger
	b, err := json.Marshal(trade)
	if err == nil {
		err = stub.PutState(trade.TradeID,b)
	} else {
		return errors.New("Error while updating trade status")
	}
	return nil
	*/
}

func (t *SimpleChaincode) trial(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, errors.New("********* TRIAL ERROR *********")
}

/* error handling
	1. uuid return error
	2. no error returned check transactionID incremented or not
	3. maintain transaction status and check every time 

*/

/* if error 
update transaction status 
dont increment transaction number or trade number
dont include transaction in trade history
*/
// read trades of a client
func (t *SimpleChaincode) readTrades(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args)== 1 {
		// read entity state
		entitybyte,err := stub.GetState(args[0])																									
		if err != nil {
			return nil, errors.New("Error while getting entity info from ledger")
		}
		var entity Entity
		err = json.Unmarshal(entitybyte, &entity)		
		if(err != nil){
			return nil, errors.New("Error while unmarshalling entity data")
		}
		trades := make([]Trade,len(entity.TradeHistory))
		for i:=0; i<len(entity.TradeHistory); i++ {
			byteVal,err := stub.GetState(entity.TradeHistory[i])
			if err != nil {
				return nil, errors.New("Error while getting trades info from ledger")
			}
			err = json.Unmarshal(byteVal, &trades[i])	
			if err != nil {
				return nil, errors.New("Error while unmarshalling trades")
			}	
		}
		b, err := json.Marshal(trades)
		if err != nil {
			return nil, errors.New("Error while marshalling trades")
		}
		return b, nil
	}
	return nil, errors.New("Incorrect number of arguments")
}
func (t *SimpleChaincode) readIssueRequests(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  return nil, nil
  /*
	var quoteTransactions []string
	// get current Trade number
	ctidByte, err := stub.GetState("currentTransactionNum")
	if(err != nil){
		return nil, errors.New("Error while getting currentTransactionNum from ledger")
	}
	tradeNum,err := strconv.Atoi(string(ctidByte))
	if(err != nil){
		return nil, errors.New("Error while converting ctidByte to integer")
	}
	// check all trades
	for tradeNum > 1000 {
		// read trade state
		tradebyte,err := stub.GetState("trans"+strconv.Itoa(tradeNum))
		if err != nil {
			return nil, errors.New("Error while getting trade info from ledger")
		}
		var trade Transaction
		err = json.Unmarshal(tradebyte, &trade)		
		if err != nil {
			return nil, errors.New("Error while unmarshalling trade data")
		}
		// check status
		if trade.Status == "New Issue" {
			quoteTransactions = append(quoteTransactions,trade.TransactionHistory[0])
		} else if trade.Status == "Responded" { // check who has responded
			respondedFlag := false
			bytes, _ := stub.GetCallerCertificate()
			x509Cert, _ := x509.ParseCertificate(bytes)
			currentUserID := x509Cert.Subject.CommonName
			
			for i:=0; i< len(trade.TransactionHistory); i++ {
				tranbyte,err := stub.GetState(trade.TransactionHistory[i])
				if(err != nil){
					return nil, errors.New("Error while getting transaction from ledger")
				}
				var tran Transaction
				err = json.Unmarshal(tranbyte, &tran)		
				if(err != nil){
					return nil, errors.New("Error while unmarshalling tran data")
				}
				if tran.TransactionType == "Response" {
					if tran.ToUser == currentUserID {
						respondedFlag = true
						break
					}
				}
			}
			if respondedFlag == false {
				quoteTransactions = append(quoteTransactions,trade.TransactionHistory[0])
			}
		}
		tradeNum--
	}
	b, err := json.Marshal(quoteTransactions)
	return b, nil
	*/
}

func updateTransactionStatus(stub shim.ChaincodeStubInterface, transactionID string, status string) (error) {
		//Transaction
		t := Transaction{
		TransactionID: transactionID,
		Status: status,
		}
		// convert to Transaction to JSON
		b, err := json.Marshal(t)
		// write to ledger
		if err == nil {
			err = stub.PutState(t.TransactionID,b)
			if(err != nil){
				return errors.New("Error while writing Transaction to ledger")
			}
		} else {
			return errors.New("Json Marshalling error")
		}
		return nil
}
func (t *SimpleChaincode) getEntityList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var allEntities []string
	var entities []string
	// get current Trade number
	ctidByte, err := stub.GetState("entityList")
	if(err != nil){
		return nil, errors.New("Error while getting entity list from ledger")
	}
	err = json.Unmarshal(ctidByte, &allEntities)		
	if(err != nil){
		return nil, errors.New("Error while unmarshalling entity data")
	}
	// check all entities
	for i:=0; i< len(allEntities); i++ {
		// read trade state
		entityByte,err := stub.GetState(allEntities[i])
		if err != nil {
			return nil, errors.New("Error while getting entity info from ledger")
		}
		var entity Entity
		err = json.Unmarshal(entityByte, &entity)		
		if err != nil {
			return nil, errors.New("Error while unmarshalling entity data")
		}
		// check type
		if entity.EntityType == "Issuer" || entity.EntityType == "Bank" || entity.EntityType == "Investor"{
			entities = append(entities,allEntities[i])
		}
	}
	b, err := json.Marshal(entities)
	return b, nil
}
func (t *SimpleChaincode) getEntities(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var allEntities []string
	//var entities []string
	// get current Trade number
	ctidByte, err := stub.GetState("entityList")
	if(err != nil){
		return nil, errors.New("Error while getting entity list from ledger")
	}
	err = json.Unmarshal(ctidByte, &allEntities)		
	if(err != nil){
		return nil, errors.New("Error while unmarshalling entity data")
	}
	// check all entities
	entities := make([]Entity,len(allEntities))

	for i:=0; i< len(allEntities); i++ {
		// read trade state
		entityByte,err := stub.GetState(allEntities[i])
		if err != nil {
			return nil, errors.New("Error while getting entity info from ledger")
		}
		//var entity Entity
		err = json.Unmarshal(entityByte, &entities[i])		
		if err != nil {
			return nil, errors.New("Error while unmarshalling entity data")
		}
		
	}
	b, err := json.Marshal(entities)
		if err != nil {
			return nil, errors.New("Error while marshalling entities")
		}
	return b, nil
}
func (t *SimpleChaincode) getAllTrades(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// check entity type
	entitybyte,err := stub.GetState(args[0])																									
	if err != nil {
		return nil, errors.New("Error while getting entity info from ledger")
	}
	var entity Entity
	err = json.Unmarshal(entitybyte, &entity)		
	if err != nil {
		return nil, errors.New("Error while unmarshalling entity data")
	}
	if entity.EntityType == "RegBody" {		
			var tradeList []string
			// get current Trade number
			ctidByte, err := stub.GetState("currentTradeNum")
			if err != nil {
				return nil, errors.New("Error while getting currentTradeNum from ledger")
			}
			tradeNum,err := strconv.Atoi(string(ctidByte))
			if err != nil {
				return nil, errors.New("Error while converting ctidByte to integer")
			}
			for tradeNum > 1000 {
					tradeList = append(tradeList,"trade"+strconv.Itoa(tradeNum))
					tradeNum--
			}
			trades := make([]Trade,len(tradeList))
			for i:=0; i<len(tradeList); i++ {
				byteVal,err := stub.GetState(tradeList[i])
				if err != nil {
					return nil, errors.New("Error while getting trades info from ledger")
				}
				err = json.Unmarshal(byteVal, &trades[i])	
				if err != nil {
					return nil, errors.New("Error while unmarshalling trades")
				}
			}
			b, err := json.Marshal(trades)
			if err != nil {
				return nil, errors.New("Error while marshalling trades")
			}
			return b, nil
	} 
	return nil, errors.New("Error only Regulatory Body can access all trades")
}
func (t *SimpleChaincode) getTransactionStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
		if len(args)== 1 {
				transactionID := "trans"+args[0]
				tbyte,err := stub.GetState(transactionID)
				if err != nil {
					return []byte("Error while getting transaction from ledger to get transaction status of "+transactionID), nil
				}
				var transaction Transaction
				err = json.Unmarshal(tbyte, &transaction)
				if err != nil {
					return []byte("Error while unmarshalling transaction data to get transaction status of "+transactionID), nil
				}
				return []byte(transaction.Status),nil
		}
		return nil, errors.New("Incorrect number of arguments")
}

// User by Issuer to Create new Issue in the Ledger
/*			arg 0 	: login user id
			arg 1	:	IOI ID
			arg 2	:	Coupon
			arg 3 	:	Rate
			arg 4	:	Price
			arg 5	:	Maturity date
			arg	6	:	Issue Date
			arg 7	:	Callable

*/
func (t *SimpleChaincode) createIssue(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//Need all parameters for the Bond Instrument
	if len(args)== 8{
		caller := args[0]
		// Check if the IOI Id already exists
		Ioibyte, err := stub.GetState(args[1])
		if err != nil {
			return nil, errors.New("Instrument with this ID already Exists, Try a different Name")
		}
		var vioi Ioi
		err = json.Unmarshal(Ioibyte, &vioi)
		if err != nil{
			return nil, errors.New("Error while unmarshalling IOI record")
		}
		
		p,err := strconv.ParseFloat(args[4],64)  // Price
		if err != nil {
			return nil,errors.New( "Error while converting Price to integer")
			
		}
		issuer := vioi.Owner
		quantity := vioi.Notional/p


		r,err := strconv.ParseFloat(args[3],64)  // Rate
		if err != nil {
			return nil,errors.New( "Error while converting Rate to integer")
		
		}

		ctidByte1,err1 := stub.GetState("currentInstrumentNum")
		if(err1 != nil){
			return nil, errors.New("Error while getting currentInstrumentNum from ledger")
		}
		instid,err := strconv.Atoi(string(ctidByte1))
		if(err != nil){
			return nil, errors.New("Error while converting ctidByte to integer")
		}
		instid = instid + 1
		instrumentID := "INST"+strconv.Itoa(instid)
		
		// convert to Instrument to JSON
		inst := Instrument {
		Symbol :instrumentID,
		Coupon :args[2],
		Quantity :int(quantity),
		InstrumentPrice :p,
		Rate :r,
		SettlementDate :args[5],
		IssueDate	:args[6],
		Callable	:args[7],
		Status :"PublishedToBank",
		Owner : caller,
		Issuer : vioi.Owner,
		}
		
		b, err := json.Marshal(inst)
		// write to ledger
		if err == nil {
			err = stub.PutState(inst.Symbol,b)
			if err != nil {
				 return nil, errors.New("Error while create new Issue")
				
			}
		} 
		
		// add Symbol ID to entity's Instrument List
		err = updateInstrumentHistory(stub, caller,inst.Symbol)
		if err != nil {
			return nil, errors.New( "Error while updating Instrument History : Caller : "+caller+" :"+inst.Symbol)
		}	
		// add Symbol ID to entity's Instrument List
		err = updateInstrumentHistory(stub, issuer,inst.Symbol)
		if err != nil {
			return nil, errors.New( "Error while updating Instrument History : Caller : "+caller+" :"+inst.Symbol)
		}
		
		ctidByte1,err1 = stub.GetState("currentTransactionNum")
		if(err1 != nil){
			return nil, errors.New("Error while getting currentTransactionNum from ledger")
		}
		tid,err := strconv.Atoi(string(ctidByte1))
		if(err != nil){
			return nil, errors.New("Error while converting ctidByte to integer")
		}
		tid = tid + 1
		transactionID := "trans"+strconv.Itoa(tid)
		
		tr := Transaction {
		TransactionID: transactionID,
		TransactionType: "Create Instrument",
		FromUser:	caller,														// 
		ToUser: issuer,  
		//SettlementDate: time.Date(year, month, day, 0, 0, 0, 0, time.UTC),				// based on input
		Status: "Success",
		Symbol:instrumentID,
		TimeStamp : time.Now().Format("2006-01-02 15:04:05"),
		}

		// convert to JSON
		b, err = json.Marshal(tr)
		
		// write to ledger
		if err == nil {
			err = stub.PutState(tr.TransactionID,b)
			if err != nil {
				_ = updateTransactionStatus(stub, transactionID, "Error while writing Response transaction to ledger")
				return nil, nil
			}
		} else {
			_ = updateTransactionStatus(stub, transactionID, "Error while marshalling transaction data")
			return nil, nil
		}
		
		err = stub.PutState("currentInstrumentNum", []byte(strconv.Itoa(instid)))
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while writing current IOI Number to ledger")
			return nil, nil
		}
		err = stub.PutState("currentTransactionNum", []byte(strconv.Itoa(tid)))
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while writing current Transaction Number to ledger")
			return nil, nil
		}
		// add Transaction ID to entity's trade history
		err = updateTradeHistory(stub, tr.ToUser, tr.TransactionID)
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while updating trade history")
			return nil, nil
		}	
		// add Transaction ID to entity's trade history
		err = updateTradeHistory(stub, tr.FromUser, tr.TransactionID)
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while updating trade history for Issuer")
			return nil, nil
		}
		vioi.Status="Responded"
		vioi.Symbol=instrumentID
		vioi.TransactionHistory=append(vioi.TransactionHistory, transactionID)
		// convert to JSON
		b, err = json.Marshal(vioi)
		
		// write to ledger
		if err == nil {
			err = stub.PutState(vioi.IoiId,b)
			if err != nil {
				_ = updateTransactionStatus(stub, transactionID, "Error while writing Response transaction to ledger")
				return nil, nil
			}
		} else {
			_ = updateTransactionStatus(stub, transactionID, "Error while marshalling IOI data")
			return nil, nil
		}

		commission := float64(quantity)*inst.InstrumentPrice*.001
		price := float64(quantity)*inst.InstrumentPrice
		err = t.updateEntityBalance(stub,caller, -1*(price - commission))   //Caller
		if err != nil {
				return nil, errors.New(err.Error())
		}
		err = t.updateEntityBalance(stub,issuer, (price - commission))   //Client
		if err != nil {
				return nil, errors.New(err.Error())
		}
		err =t.updateInstrumentTradeHistory(stub, inst.Symbol, transactionID)
		if err != nil {
				return nil, errors.New(err.Error())
		}

		return []byte(inst.Symbol), nil
	}
	return nil, errors.New("Incorrect number of arguments")
}

func (t *SimpleChaincode) getInstrument(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
		instbyte,err := stub.GetState(args[0])																									
		if err != nil {
			return nil, errors.New("Error while getting Instrument info from ledger")
		}
		return instbyte, nil
}
func (t *SimpleChaincode) updateInstrumentStatus(stub shim.ChaincodeStubInterface, symbol string, possassion string, status string) (error) {
		instbyte,err := stub.GetState(symbol)																									
		if err != nil {
			return  errors.New("Error while getting Instrument info from ledger")
		}
		var inst Instrument
		err = json.Unmarshal(instbyte, &inst)
		if err != nil {
			return  errors.New("Unable to Unmarshal Instrument")
		}
		entityByte , err := stub.GetState(possassion)
		if err != nil {
		   return errors.New("Error in getting Entity")
		}
		var entity Entity
		err = json.Unmarshal(entityByte, &entity)
		if err != nil {
		   return  errors.New("Error in getting Entity Unarshal")
		}
		if entity.EntityType =="Bank" {
			inst.Bank = possassion
		}
		
		inst.Owner = possassion
		inst.Status = status
		b , err := json.Marshal(inst)
		if err != nil {
			return  errors.New("Unable to marshal Instrument")
		}
		err = stub.PutState(inst.Symbol,b)
		if err != nil {
			return  errors.New("Unable to update Instrument status")
		}
		return  nil
}


func (t *SimpleChaincode) updateInstrumentTradeHistory(stub shim.ChaincodeStubInterface, symbol string, TransactionID string) (error) {
		instbyte,err := stub.GetState(symbol)																									
		if err != nil {
			return  errors.New("Error while getting Instrument info from ledger")
		}
		var inst Instrument
		err = json.Unmarshal(instbyte, &inst)
		if err != nil {
			return  errors.New("Unable to Unmarshal Instrument")
		}
		newTradeList := []string{}
		
		for _, v := range inst.TradeID {
			if v != TransactionID {
			newTradeList = append(newTradeList, v)
			}
		}
		newTradeList = append(newTradeList,TransactionID)
		inst.TradeID = newTradeList
		b , err := json.Marshal(inst)
		if err != nil {
			return  errors.New("Unable to marshal Instrument")
		}
		err = stub.PutState(inst.Symbol,b)
		if err != nil {
			return  errors.New("Unable to update Instrument Trade ")
		}
		return  nil
}

func (t *SimpleChaincode) getAllInstruments(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
/*
	args 0 : Entity
	args 1 : Status
*/
	// check entity type
	entitybyte,err := stub.GetState(args[0])																									
	if err != nil {
		return nil, errors.New("Error while getting entity info from ledger")
	}
	var entity Entity
	err = json.Unmarshal(entitybyte, &entity)		
	if err != nil {
		return nil, errors.New("Error while unmarshalling entity data")
	}
	fmt.Println("Status and Entity" + args[1])
	status := args[1]
	//if entity.EntityType == "RegBody" {		
	instruments := make([]Instrument,len(entity.Instruments))
	var instrumentArray []Instrument //:= make([]Instrument,1)

		for i:=0; i<len(entity.Instruments); i++ {
			byteVal,err := stub.GetState(entity.Instruments[i])
			if err != nil {
				return nil, errors.New("Error while getting Instrument info from ledger")
			}
			fmt.Println("Bytevalue of Instruent and Entity" + string(byteVal))
			
			err = json.Unmarshal(byteVal, &instruments[i])	
			if err != nil {
				return nil, errors.New("Error while unmarshalling trades")
			}
			if entity.EntityType =="Issuer" && status =="Outstanding" && instruments[i].Status != "New Issue"{
			
				instrumentArray = append(instrumentArray,instruments[i])
				
			}else if entity.EntityType =="Bank" && status =="Outstanding" && instruments[i].Status != "PublishToBank"{
					instrumentArray = append(instrumentArray,instruments[i])
					
			}else if entity.EntityType =="Investor" && status =="Outstanding" && instruments[i].Status != "PublishToInvestor"{
					instrumentArray = append(instrumentArray,instruments[i])
					
			}else if instruments[i].Status == status{
			
				instrumentArray = append(instrumentArray,instruments[i])
			
			}else {
			 fmt.Println("status" + status + " -" +instruments[i].Status )
			}
		}
		
		if status == "All" {
			instrumentArray = instruments
		}
		b, err := json.Marshal(instrumentArray)
		if err != nil {
			return nil, errors.New("Error while marshalling trades")
		}
		return b, nil

}

func (t *SimpleChaincode) getAllInstrumentTrades(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
/*
	args 0 : Entity
	args 1 : Symbol

*/
	// check entity type
	entitybyte,err := stub.GetState(args[0])																									
	if err != nil {
		return nil, errors.New("Error while getting entity info from ledger")
	}
	var entity Entity
	err = json.Unmarshal(entitybyte, &entity)		
	if err != nil {
		return nil, errors.New("Error while unmarshalling entity data")
	}
	
	instbyte,err := stub.GetState(args[1])																									
	if err != nil {
		return nil, errors.New("Error while getting Instrument info from ledger")
	}
	var instrument Instrument
	fmt.Println("Instruent Byte "+string(instbyte))
	err = json.Unmarshal(instbyte, &instrument)		
	if err != nil {
		return nil, errors.New("Error while unmarshalling Instrument data")
	}
	
	trades := make([]Transaction,len(instrument.TradeID))
	var tradesArray  []Transaction
	
		for i:=0; i<len(instrument.TradeID); i++ {
			byteVal,err := stub.GetState(instrument.TradeID[i])
			if err != nil {
				return nil, errors.New("Error while getting Transaction info from ledger")
			}
			err = json.Unmarshal(byteVal, &trades[i])	
			if err != nil {
				return nil, errors.New("Error while unmarshalling trades")
			}
			if (trades[i].FromUser == args[0] ||trades[i].ToUser == args[0]){
			tradesArray = append(tradesArray,trades[i])

			}
		}
		
		fmt.Println("Entity Type "+entity.EntityType)
		if entity.EntityType =="RegBody"{
			tradesArray = trades
		}
		b, err := json.Marshal(tradesArray)
		fmt.Println("Trades List "+string(b))
		
		if err != nil {
			return nil, errors.New("Error while marshalling Transacations")
		}
		return b, nil

}


func (t *SimpleChaincode) updateEntityBalance(stub shim.ChaincodeStubInterface, entity string, Amount float64) (error) {

// check entity type
	entitybyte,err := stub.GetState(entity)																									
	if err != nil {
		return errors.New("Error while getting entity info from ledger")
	}
	var entity1 Entity
	err = json.Unmarshal(entitybyte, &entity1)		
	if err != nil {
		return  errors.New("Error while unmarshalling entity data")
	}
	
	entity1.Balance = entity1.Balance+ Amount
	if entity1.Balance <0 {
	  return  errors.New("Inssufficient Balance for Entity" + entity)
	}

	b, err := json.Marshal(entity1)
	if err != nil {
		return  errors.New("Error while marshal entity data")
	}
	err = stub.PutState(entity, b)
	if err != nil {
		return  errors.New("Error while updating entity data")
	}
	return  nil
}

func (t *SimpleChaincode) payCoupon(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
			/*
				args 0 : Symbol
			*/
		instbyte , err := stub.GetState(args[0])
		if err != nil {
			 return nil, errors.New("Error in finding instrument")
		}
		var inst Instrument
		err = json.Unmarshal(instbyte, &inst)		
		if err != nil {
			return  nil,errors.New("Error while unmarshalling instbyte data")
		}
		
		owner  :=  inst.Owner
		issuer  := inst.Issuer
		coupon  := inst.Rate * float64(inst.Quantity)/100
		
		entitybyte,err := stub.GetState(issuer)																									
		if err != nil {
			return nil,errors.New("Error while getting Issuer info from ledger")
		}
		var entity Entity
		err = json.Unmarshal(entitybyte, &entity)		
		if err != nil {
			return  nil,errors.New("Error while unmarshalling entity data")
		}
		
		entity.Balance = entity.Balance- coupon
		if entity.Balance <0 {
		  return  nil,errors.New("Inssufficient Balance for Entity" + issuer)
		}

		b, err := json.Marshal(entity)
		if err != nil {
			return  nil,errors.New("Error while marshal entity data")
		}
		err = stub.PutState(issuer, b)
		if err != nil {
			return  nil,errors.New("Error while updating entity data")
		}
		
		entitybyte,err = stub.GetState(owner)																									
		if err != nil {
			return nil,errors.New("Error while getting Banker/Investor info from ledger")
		}

		err = json.Unmarshal(entitybyte, &entity)		
		if err != nil {
			return  nil,errors.New("Error while unmarshalling entity data")
		}
		
		entity.Balance = entity.Balance+ coupon

		b, err = json.Marshal(entity)
		if err != nil {
			return  nil,errors.New("Error while marshal Investor/Banker entity data")
		}
		err = stub.PutState(owner, b)
		if err != nil {
			return  nil,errors.New("Error while updating entity data")
		}
		return  nil,nil
}


func (t *SimpleChaincode) issueCallout(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
			/*
				args 0 : Symbol
			*/
		instbyte , err := stub.GetState(args[0])
		if err != nil {
			 return nil, errors.New("Error in finding instrument")
		}
		var inst Instrument
		err = json.Unmarshal(instbyte, &inst)		
		if err != nil {
			return  nil,errors.New("Error while unmarshalling instbyte data")
		}
		
		owner  :=  inst.Owner
		issuer  := inst.Issuer
		price  := inst.InstrumentPrice * float64(inst.Quantity)
		inst.Status = "Expired"
		inst.Owner = issuer
		
		entitybyte,err := stub.GetState(issuer)																									
		if err != nil {
			return nil,errors.New("Error while getting Issuer info from ledger")
		}
		var entity Entity
		err = json.Unmarshal(entitybyte, &entity)		
		if err != nil {
			return  nil,errors.New("Error while unmarshalling entity data")
		}
		
		entity.Balance = entity.Balance- price
		if entity.Balance <0 {
		  return  nil,errors.New("Inssufficient Balance for Entity :" + issuer)
		}

		b, err := json.Marshal(entity)
		if err != nil {
			return  nil,errors.New("Error while marshal entity data")
		}
		err = stub.PutState(issuer, b)
		if err != nil {
			return  nil,errors.New("Error while updating entity data")
		}
		
		entitybyte,err = stub.GetState(owner)																									
		if err != nil {
			return nil,errors.New("Error while getting Banker/Investor info from ledger")
		}

		err = json.Unmarshal(entitybyte, &entity)		
		if err != nil {
			return  nil,errors.New("Error while unmarshalling entity data")
		}
		
		entity.Balance = entity.Balance+ price

		b, err = json.Marshal(entity)
		if err != nil {
			return  nil,errors.New("Error while marshal Investor/Banker entity data")
		}
		err = stub.PutState(owner, b)
		if err != nil {
			return  nil,errors.New("Error while updating entity data")
		}
		
		b , err = json.Marshal(inst)
		if err != nil {
			return  nil,errors.New("Error while marshal Instrument data")
		}
		err = stub.PutState(args[0], b)
		if err != nil {
			return  nil,errors.New("Error while updating entity data")
		}
		return  nil,nil
}

/*  RequestForInstrument
		args 0: Calling User id
		args 1: Bank Id
		args 2: Notional
		args 3: Tenor
*/
func (t *SimpleChaincode) requestForInstrument(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args)== 4 {
		caller := args[0]
		bank := args[1]
		notional, err := strconv.ParseFloat(args[2], 64)
		tenor := args[3]
		
		quoteByte, err := stub.GetState("currentIoiNum")
		if err != nil {
			return nil,errors.New("Unable to find current Rfq number")
		}
		
		qtid,err := strconv.Atoi(string(quoteByte))
		if(err != nil){
			return nil, errors.New("Error while converting quoteByte to integer")
		}
		qtid = qtid + 1
		IoiID := "IOI"+strconv.Itoa(qtid)
		
		ctidByte1,err1 := stub.GetState("currentTransactionNum")
		if(err1 != nil){
			return nil, errors.New("Error while getting currentTransactionNum from ledger")
		}
		tid,err := strconv.Atoi(string(ctidByte1))
		if(err != nil){
			return nil, errors.New("Error while converting ctidByte to integer")
		}
		tid = tid + 1
		transactionID := "trans"+strconv.Itoa(tid)
		
		
		// get bank's enrollment id
		bytes, err := stub.GetCallerCertificate();
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while getting caller certificate")
			return nil, nil
		}
		x509Cert, err := x509.ParseCertificate(bytes);
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while parsing caller certificate")
			return nil, nil
		}		
		fmt.Println("Create IOI : x509Cert"+x509Cert.Subject.CommonName)
		
		ioi := Ioi {
		IoiId:IoiID, 				// ioi/rfq id
		Notional:notional, 
		Tenor:tenor,
		Bank:bank,
		Status:"New IOI", 				// "New" or "Responded" 
		Owner:caller,

		}
		
		ioi.TransactionHistory = append(ioi.TransactionHistory,transactionID)
		d, err := json.Marshal(ioi)

		// write to ledger
		if err == nil {
			err = stub.PutState(ioi.IoiId,d)
			if err != nil {
				return nil, errors.New("Error while writing Response transaction to ledger")
			}
		} else {
			return nil, errors.New("Error while unMarshalling Response Ioi to ledger")

		}
		
		tr := Transaction {
		TransactionID: transactionID,
		TransactionType: "New IOI",
		FromUser:	caller,														// 
		ToUser: bank,  //x509Cert.Subject.CommonName,											// 
		//SettlementDate: time.Date(year, month, day, 0, 0, 0, 0, time.UTC),				// based on input
		Status: "Success",
		Symbol:IoiID,
		TimeStamp : time.Now().Format("2006-01-02 15:04:05"),
		}

		// convert to JSON
		b, err := json.Marshal(tr)
		
		// write to ledger
		if err == nil {
			err = stub.PutState(tr.TransactionID,b)
			if err != nil {
				_ = updateTransactionStatus(stub, transactionID, "Error while writing Response transaction to ledger")
				return nil, nil
			}
		} else {
			_ = updateTransactionStatus(stub, transactionID, "Error while marshalling transaction data")
			return nil, nil
		}
		
		err = stub.PutState("currentIoiNum", []byte(strconv.Itoa(qtid)))
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while writing current IOI Number to ledger")
			return nil, nil
		}
		err = stub.PutState("currentTransactionNum", []byte(strconv.Itoa(tid)))
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while writing current Transaction Number to ledger")
			return nil, nil
		}
		
		
		// add Transaction ID to entity's trade history
		err = updateTradeHistory(stub, tr.ToUser, tr.TransactionID)
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while updating trade history")
			return nil, nil
		}	
		// add Transaction ID to entity's trade history
		err = updateTradeHistory(stub, tr.FromUser, tr.TransactionID)
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while updating trade history for Issuer")
			return nil, nil
		}
		// add IOI ID to entity's  history
		err = updateIOIHistory(stub, tr.ToUser, IoiID)
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while updating trade history")
			return nil, nil
		}	
		// add IOI ID to entity's  history
		err = updateIOIHistory(stub, tr.FromUser, IoiID)
		if err != nil {
			_ = updateTransactionStatus(stub, transactionID, "Error while updating trade history")
			return nil, nil
		}	
		
		return nil, nil
	}
	return nil, errors.New("Incorrect number of arguments")
}

func updateIOIHistory(stub shim.ChaincodeStubInterface, entityID string, IOIID string) (error) {
	// read entity state
	entitybyte,err := stub.GetState(entityID)																										
	if err != nil {
		return errors.New("Error while getting entity info from ledger")
	}
	var entity Entity
	err = json.Unmarshal(entitybyte, &entity)		
	if err != nil {
		return errors.New("Error while unmarshalling entity data")
	}
	// add IOIID to history
	entity.IoiList = append(entity.IoiList,IOIID)
	// write entity state to ledger
	b, err := json.Marshal(entity)
	if err == nil {
		err = stub.PutState(entity.EntityID,b)
	} else {
		return errors.New("Error while updating entity status")
	}
	return nil
}

func (t *SimpleChaincode) getAllIoi(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
/*
	args 0 : Entity
*/
	// check entity type
	entitybyte,err := stub.GetState(args[0])																									
	if err != nil {
		return nil, errors.New("Error while getting entity info from ledger")
	}
	var entity Entity
	err = json.Unmarshal(entitybyte, &entity)		
	if err != nil {
		return nil, errors.New("Error while unmarshalling entity data")
	}

	ioi := make([]Ioi,len(entity.IoiList))
	var IoiArray []Ioi //:= make([]Ioi,1)

		for i:=0; i<len(entity.IoiList); i++ {
			byteVal,err := stub.GetState(entity.IoiList[i])
			if err != nil {
				return nil, errors.New("Error while getting Instrument info from ledger")
			}
			fmt.Println("Bytevalue of Instruent and Entity" + string(byteVal))
			
			err = json.Unmarshal(byteVal, &ioi[i])	
			if err != nil {
				return nil, errors.New("Error while unmarshalling ioi")
			}

		}
		
		IoiArray = ioi

		b, err := json.Marshal(IoiArray)
		if err != nil {
			return nil, errors.New("Error while marshalling trades")
		}
		return b, nil

}
