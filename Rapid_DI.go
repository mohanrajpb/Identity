package main

import (
	"errors"
	"fmt"

	//"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("mylogger")

type IdentityChain struct {
}

//custom data models

type ItemPID struct {
  ItemType string `json:"itemtype"`
  PIDName   string `json:"pidname"`
  PUUID     string `json:"puuid"`
  PIDStatus string `json:"pidstatus"`
}


type ItemSIDs struct {
	SIDName            string `json:"sidname"`
	SID                string `json:"sid"`
  SIDDateIncluded    string `json:"siddateincluded"`
	SIDValidFrom       string `json:"sidvalidfrom"`
  SIDValidTill       string `json:sidvalidtill"`
  SIDValidationType  string `json:"sidvalidationtype"`
  SIDAuthorisedBy    string `json:"sidauthorisedby"`
}

type DocketAccess struct {
	DAIDAccessed         string `json:"daidaccessed"`
	DAReceiver          string `json:"dareceiver"`
  DADateApplied       string `json:"dadateapplied"`
  DADateAuthorised    string `json:"dadateauthorised"`
	DAValidFrom         string `json:"davalidfrom"`
  DAValidTill         string `json:"davalidtill"`
}

type Docket struct {
  ItemPID        ItemPID    `json:"itempid"`
  ItemSIDs       ItemSIDs  `json:"itemsids"`
  DocketAccess   DocketAccess `json:"docketaccess"`
}

func GetDocket(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering GetDocket")

	if len(args) < 1 {
		logger.Error("Invalid number of arguments")
		return nil, errors.New("Missing Item UUID")
	}

	var itemUUID = args[0]
	bytes, err := stub.GetState(itemUUID)
	if err != nil {
		logger.Error("Could not fetch Docket with identifier "+itemUUID+" from ledger", err)
		return nil, err
	}
	return bytes, nil
}

func CreateDocket(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering CreateDocket")

	if len(args) < 2 {
		logger.Error("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for loan application creation")
	}

	var itemUUID = args[0]
	var docketUpdate = args[1]

	err := stub.PutState(itemUUID, []byte(docketUpdate))
	if err != nil {
		logger.Error("Could not save identity or docket to ledger", err)
		return nil, err
	}

	var customEvent = "{eventType: 'docketCreation', description:" + itemUUID + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully created identifier and docket")
	return nil, nil

}

/**
Updates the status of the loan application

func UpdateDocket(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering UpdateDocket")

	if len(args) < 2 {
		logger.Error("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for loan application update")
	}

	var itemUUID = args[0]
	var status = args[1]

	laBytes, err := stub.GetState(loanApplicationId)
	if err != nil {
		logger.Error("Could not fetch loan application from ledger", err)
		return nil, err
	}
	var loanApplication LoanApplication
	err = json.Unmarshal(laBytes, &loanApplication)
	loanApplication.Status = status

	laBytes, err = json.Marshal(&loanApplication)
	if err != nil {
		logger.Error("Could not marshal loan application post update", err)
		return nil, err
	}

	err = stub.PutState(loanApplicationId, laBytes)
	if err != nil {
		logger.Error("Could not save loan application post update", err)
		return nil, err
	}

	var customEvent = "{eventType: 'loanApplicationUpdate', description:" + loanApplicationId + "' Successfully updated status'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully updated loan application")
	return nil, nil

}
**/
func (t *IdentityChain) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

func (t *IdentityChain) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "GetDocket" {
		return GetDocket(stub, args)
	}
	return nil, nil
}

func GetCertAttribute(stub shim.ChaincodeStubInterface, attributeName string) (string, error) {
	logger.Debug("Entering GetCertAttribute")
	attr, err := stub.ReadCertAttribute(attributeName)
	if err != nil {
		return "", errors.New("Couldn't get attribute " + attributeName + ". Error: " + err.Error())
	}
	attrString := string(attr)
	return attrString, nil
}

func (t *IdentityChain) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "CreateDocket" {
		username, _ := GetCertAttribute(stub, "username")
		role, _ := GetCertAttribute(stub, "role")
		if role == "Bank_Home_Loan_Admin" {
			return CreateDocket(stub, args)
		} else {
			return nil, errors.New(username + " with role " + role + " does not have access to create a loan application")
		}

	}
	return nil, nil
}

type customEvent struct {
	Type       string `json:"type"`
	Decription string `json:"description"`
}

func main() {

	lld, _ := shim.LogLevel("DEBUG")
	fmt.Println(lld)

	logger.SetLevel(lld)
	fmt.Println(logger.IsEnabledFor(lld))

	err := shim.Start(new(IdentityChain))
	if err != nil {
		logger.Error("Could not start IdentityChain")
	} else {
		logger.Info("IdentityChain successfully started")
	}
}
