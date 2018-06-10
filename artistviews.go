package main

import (
	"fmt"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type ArtistViewsChaincode struct {
}

//var A, B string
//var Aval, Bval, X int

// Init callback representing the invocation of a chaincode
// This chaincode will manage two accounts A and B and will transfer X units from A to B upon invoke
func (t *ArtistViewsChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Init Called")
	return shim.Success(nil)
}

func (t *ArtistViewsChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Invoke Called")
	return shim.Success(nill)
}

func (t *ArtistViewsChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	
	fmt.Println("Invoke Called")
	fmt.Println("Hello world")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		return t.invoke(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\"")
}

func main() {
	err := shim.Start(new(ArtistViewsChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
