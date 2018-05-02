package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Pet struct {
	Name string `json:"name"`
	Picture string `json:"picture"`
	Breed string `json:"breed"`
	Location string `json:"location"`
	Age int `json:"age"`
	Owner string `json:"owner"`
}

/*
 * The Init method is called when the Smart Contract "pet-shop" is instantiated by the blockchain network
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "pet-shop"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryAllPets" {
		return s.queryAllPets(APIstub)
	} else if function == "adoptPet" {
		return s.adoptPet(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}


/*
* The initLedger method *
Will add test data to our network
*/
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	pets := []Pet{
		Pet{Name: "Frieda", Picture: "images/scottish-terrier.jpeg", Age: 3, Breed: "Scottish Terrier", Location: "Lisco, Alabama"},
		Pet{Name: "Gina", Picture: "images/scottish-terrier.jpeg", Age: 3,Breed: "Scottish Terrier", Location: "Tooleville, West Virginia"},
		Pet{Name: "Collins", Picture: "images/french-bulldog.jpeg", Age: 2,Breed: "French Bulldog", Location: "Freeburn, Idaho"},
		Pet{Name: "Melissa", Picture: "images/boxer.jpeg", Age: 2,Breed: "Boxer", Location: "Camas, Pennsylvania"},
	}

	i := 0
	for i < len(pets) {
		petAsBytes, _ := json.Marshal(pets[i])
		APIstub.PutState(strconv.Itoa(i+1), petAsBytes)
		i = i + 1
	}

	return shim.Success(nil)
}


/*
* The queryAllPets method *
allows for assessing all the records added to the ledger(all pets)
This method does not take any arguments. Returns JSON string containing results. 
*/
func (s *SmartContract) queryAllPets(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllPets:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
* The adoptPet method *
allows the user to adopt a pet
This function takes in 1 argument (pet id). 
*/
func (s *SmartContract) adoptPet(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	petAsBytes, _ := APIstub.GetState(args[0])
	if petAsBytes == nil {
		return shim.Error("Could not locate pet")
	}
	pet := Pet{}

	json.Unmarshal(petAsBytes, &pet)

	// GetCreator returns marshaled serialized identity of the client
	creator, err := APIstub.GetCreator()
	if err != nil {
		return shim.Error(fmt.Sprintf("Error received on GetCreator", err))
	}
	certStart := bytes.IndexAny(creator, "----BEGIN CERTIFICATE-----")
	if certStart == -1 {
		return shim.Error("No certificate found")
	}
	certText := creator[certStart:]
	block, _ := pem.Decode(certText)
	if block == nil {
		return shim.Error(fmt.Sprintf("Error received on pem.Decode of certificate",  certText))
	}
	ucert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error received on ParseCertificate", err))
	}
	pet.Owner = ucert.Subject.CommonName

	petAsBytes, _ = json.Marshal(pet)
	err = APIstub.PutState(args[0], petAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change pet owner: %s", args[0]))
	}

	return shim.Success(nil)
}


/*
* main function *
calls the Start function 
The main function starts the chaincode in the container during instantiation.
*/
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
